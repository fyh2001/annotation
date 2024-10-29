<script setup>
import {
  NearMeRound,
  SquareRound,
  Brightness1Round,
  CoPresentSharp,
} from "@vicons/material";
import { ref, onMounted, onUnmounted } from "vue";
import { getTask, submit } from "../../api/methods/task";
import { baseURL } from "../../api";

const getNewTaskBtnDisable = ref(false);
const currentTaskId = ref(null);

const getNewTask = async () => {
  const { data } = await getTask();

  if (data.code !== 200) return;

  const task = data.data;

  localStorage.setItem("task", JSON.stringify(task));
  loadImage(baseURL.substring(0, baseURL.length - 3) + task.url);
  getNewTaskBtnDisable.value = true;

  currentTaskId.value = task.id;
};

const handleSubmit = async (mode) => {
  if (boxes.value.length === 0) return;

  const task = JSON.parse(localStorage.getItem("task"));
  if (new Date(task.expiredAt) < Date.now()) {
    reset();
    if (mode == "exit") {
      return;
    } else if (mode == "continue") {
      return getNewTask();
    }
  }

  const submitData = boxes.value.map((item) => {
    const rectX = item.x - imageConfig.value.x;
    const rectY = item.y - imageConfig.value.y;

    const normalizedX = rectX / imageConfig.value.width;
    const normalizedY = rectY / imageConfig.value.height;
    const normalizedWidth = item.width / imageConfig.value.width;
    const normalizedHeight = item.height / imageConfig.value.height;

    return {
      classification: item.classification,
      x: normalizedX + normalizedWidth / 2,
      y: normalizedY + normalizedHeight / 2,
      width: normalizedWidth,
      height: normalizedHeight,
    };
  });

  task.result = JSON.stringify(submitData);

  const { data } = await submit(task);

  if (data.code === 200) {
    reset();
    if (mode == "exit") {
      return;
    } else if (mode == "continue") {
      return getNewTask();
    }
  }
};

const stageConfig = ref({
  width: window.innerWidth * 0.6,
  height: window.innerHeight * 0.8,
  scaleX: 1,
  scaleY: 1,
});

const stageRef = ref(null);
const transformerRef = ref(null);

const image = ref(null);
const imageWidth = ref(0);
const imageHeight = ref(0);
const imageConfig = ref(null);

const boxes = ref([]);
const scale = ref(1);
const selectedRectIndex = ref(null);
const tempSelectedRectIndex = ref(null);
const isDrawing = ref(false);
const undoStack = ref([]);
const currentMode = ref("rectTool");

const currentClassification = ref(0);
const classifications = ref({
  0: {
    id: 0,
    name: "分类1",
    color: "#FF0000",
  },
});

const reset = () => {
  currentTaskId.value = null;

  getNewTaskBtnDisable.value = false;
  transformerRef.value = null;

  stageConfig.value = {
    width: window.innerWidth * 0.6,
    height: window.innerHeight * 0.8,
    scaleX: 1,
    scaleY: 1,
  };

  const stage = stageRef.value.getStage();
  stage.scale({ x: 1, y: 1 });
  stage.position({ x: 0, y: 0 });
  stage.width(window.innerWidth * 0.6);
  stage.height(window.innerHeight * 0.8);

  image.value = null;
  imageConfig.value = {
    ...imageConfig.value,
    width: 0,
    height: 0,
    img: null,
  };

  boxes.value = [];
  scale.value = 1;
  selectedRectIndex.value = null;
  tempSelectedRectIndex.value = null;
  isDrawing.value = false;
  undoStack.value = [];
  currentMode.value = "rectTool";

  stage.batchDraw();
};

const selectedRect = computed(() => {
  const index = selectedRectIndex.value;

  if (index !== null && index >= 0 && index < boxes.value.length) {
    const stage = stageRef.value.getStage();
    const rect = stage.find("Rect")[index];

    if (rect && rect.x !== undefined && rect.y !== undefined) return rect;
  }

  return null;
});

const loadImage = (url) => {
  const img = new Image();
  img.src = url;
  img.onload = () => {
    image.value = img;
    const stage = stageRef.value.getStage();

    imageWidth.value = img.width;
    imageHeight.value = img.height;

    const stageWidth = stage.width();
    const stageHeight = stage.height();

    let scale = Math.min(
      stageWidth / imageWidth.value,
      stageHeight / imageHeight.value
    );

    imageConfig.value = {
      x: (stageWidth - imageWidth.value * scale) / 2,
      y: (stageHeight - imageHeight.value * scale) / 2,
      width: imageWidth.value * scale,
      height: imageHeight.value * scale,
      image: img,
    };

    stage.batchDraw();
  };
  img.onerror = (error) => {
    console.error("Image loading error:", error);
  };
};

const handleZoom = (e) => {
  e.evt.preventDefault();
  const stage = stageRef.value.getStage();
  const scaleBy = 1.05;
  const oldScale = stage.scaleX();
  const pointer = stage.getPointerPosition();

  const newScale = e.evt.deltaY > 0 ? oldScale * scaleBy : oldScale / scaleBy;
  scale.value = newScale;

  stage.scale({ x: newScale, y: newScale });

  const mousePointTo = {
    x: (pointer.x - stage.x()) / oldScale,
    y: (pointer.y - stage.y()) / oldScale,
  };

  const newPos = {
    x: pointer.x - mousePointTo.x * newScale,
    y: pointer.y - mousePointTo.y * newScale,
  };

  stage.position(newPos);
  stage.batchDraw();
};

const handleMouseDown = (e) => {
  const stage = stageRef.value.getStage();
  const pointerPos = stage.getPointerPosition();

  if (e.target === stage || currentMode.value !== "rectTool") {
    return (selectedRectIndex.value = null);
  }

  const transform = stage.getAbsoluteTransform().copy();
  transform.invert();
  const adjustedPointerPos = transform.point(pointerPos);

  boxes.value.push({
    classification: currentClassification.value,
    x: adjustedPointerPos.x,
    y: adjustedPointerPos.y,
    width: 0,
    height: 0,
    stroke: classifications.value[currentClassification.value].color,
    strokeWidth: 2,
    draggable: true,
  });

  tempSelectedRectIndex.value = boxes.value.length - 1;
  isDrawing.value = true;
};

const handleMouseMove = () => {
  if (isDrawing.value) {
    const stage = stageRef.value.getStage();
    const pointerPos = stage.getPointerPosition();
    const box = boxes.value[tempSelectedRectIndex.value];

    const transform = stage.getAbsoluteTransform().copy();
    transform.invert();

    const adjustedPointerPos = transform.point(pointerPos);

    box.width = adjustedPointerPos.x - box.x;
    box.height = adjustedPointerPos.y - box.y;

    stage.batchDraw();
  }
};

const handleMouseUp = () => {
  if (isDrawing.value) {
    isDrawing.value = false;

    undoStack.value.push([...boxes.value]);
    selectedRectIndex.value = tempSelectedRectIndex.value;
  }
};

const handleDragMove = (index, event) => {
  const box = boxes.value[index];
  box.x = event.target.x();
  box.y = event.target.y();

  const stage = stageRef.value.getStage();
  stage.batchDraw();
};

const handleDragEnd = (index, event) => {
  undoStack.value.push(JSON.parse(JSON.stringify(boxes.value)));
};

const handleRectClick = async (index) => {
  selectedRectIndex.value = index;

  const stage = stageRef.value.getStage();
  const rect = stage.find("Rect")[index];

  await nextTick();

  if (transformerRef.value && rect) {
    transformerRef.value.nodes([rect]);
    transformerRef.value.getLayer().batchDraw();
  } else {
    console.error("Transformer or selected rectangle not available.");
  }
};

const undo = () => {
  if (undoStack.value.length > 0) {
    undoStack.value.pop();
    const previousState = undoStack.value[undoStack.value.length - 1];

    if (previousState) {
      boxes.value = [...previousState];
    } else {
      boxes.value = [];
    }
    const stage = stageRef.value.getStage();
    stage.batchDraw();
  }
};

const deleteRect = () => {
  if (selectedRectIndex.value !== null) {
    boxes.value.splice(selectedRectIndex.value, 1);
    selectedRectIndex.value = null;
    undoStack.value.push([...boxes.value]);
  }
};

const handleKeydown = (event) => {
  if (event.key.toLowerCase() === "n" && mode === "selectBox") {
    boxes.value.push({
      x: 50,
      y: 50,
      width: 100,
      height: 100,
      stroke: "#FF0000",
      strokeWidth: 2,
    });
    undoStack.value.push([...boxes.value]);
  } else if (
    event.key.toLowerCase() === "delete" ||
    event.key.toLowerCase() === "backspace"
  ) {
    deleteRect();
  } else if (
    (event.ctrlKey || event.metaKey) &&
    event.key.toLowerCase() === "z"
  ) {
    undo();
  } else if (event.key.toLowerCase() === "a") {
    currentMode.value = "selectBox";
  } else if (event.key.toLowerCase() === "u") {
    currentMode.value = "rectTool";
  }
};

onMounted(() => {
  window.addEventListener("keydown", handleKeydown);

  const task = JSON.parse(localStorage.getItem("task"));

  if (task && new Date(task.expiredAt) > Date.now()) {
    getNewTaskBtnDisable.value = true;
    loadImage(baseURL.substring(0, baseURL.length - 3) + task.url);
  }
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeydown);
});
</script>

<template>
  <div>
    <div class="flex justify-between items-center px-4 py-2 border-b">
      <n-button
        type="info"
        @click="getNewTask"
        :disabled="getNewTaskBtnDisable"
      >
        获取新任务
      </n-button>
      <div>当前task ID: {{ currentTaskId }}</div>
      <div class="flex gap-4">
        <n-button strong secondary type="error" @click="handleSubmit('exit')">
          提交退出
        </n-button>
        <n-button
          strong
          secondary
          type="primary"
          @click="handleSubmit('continue')"
        >
          提交继续
        </n-button>
      </div>
    </div>
    <div class="flex justify-between py-4 px-10">
      <div class="flex flex-col gap-4 mt-8">
        <n-button
          strong
          secondary
          :type="currentMode === 'selectBox' ? 'primary' : 'tertiary'"
          @click="currentMode = 'selectBox'"
        >
          选择
          <template #icon>
            <n-icon :component="NearMeRound" size="24" />
          </template>
        </n-button>
        <n-button
          strong
          secondary
          :type="currentMode === 'rectTool' ? 'primary' : 'tertiary'"
          @click="currentMode = 'rectTool'"
        >
          矩形工具
          <template #icon>
            <n-icon :component="SquareRound" size="24" />
          </template>
        </n-button>
      </div>
      <div class="border border-2">
        <v-stage
          ref="stageRef"
          :config="stageConfig"
          @wheel="handleZoom"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
        >
          <v-layer>
            <v-image :image="image" :config="imageConfig" />
            <v-rect
              v-for="(rect, index) in boxes"
              :key="index"
              :config="rect"
              @dragmove="handleDragMove(index, $event)"
              @dragend="handleDragEnd(index, $event)"
              @click="handleRectClick(index)"
            />
            <v-transformer
              v-if="selectedRect && selectedRectIndex !== null"
              ref="transformerRef"
              :nodes="[selectedRect]"
            />
          </v-layer>
        </v-stage>
      </div>
      <div>
        <div class="p-4 mt-8 border border-1 rounded-md">
          <div>分类</div>

          <div
            class="mt-4"
            v-for="(item, index) in classifications"
            :key="index"
          >
            <n-button strong ghost type="error">
              {{ item.id }} - {{ item.name }}
              <template #icon>
                <n-icon
                  :component="Brightness1Round"
                  size="16"
                  :color="item.color"
                />
              </template>
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
body {
  margin: 0;
  overflow: hidden;
}
</style>
