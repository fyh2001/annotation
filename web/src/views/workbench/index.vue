<script setup>
import { NearMeRound, SquareRound, Brightness1Round } from "@vicons/material";
import { ref, onMounted, onUnmounted } from "vue";
import imageUrl from "../../../public/1.jpg";

// 初始化 Konva 画布配置
const stageConfig = ref({
    width: window.innerWidth * 0.6,
    height: window.innerHeight * 0.8,
    scaleX: 1,
    scaleY: 1,
});

const stageRef = ref(null);
const transformerRef = ref(null);

// 存储图像、标注框和缩放比例
const image = ref(null); // 存放加载的图像
const imageConfig = ref(null);

const boxes = ref([]); // 存放标注框
const scale = ref(1); // 缩放比例
const selectedRectIndex = ref(null); // 当前选中的矩形框索引
const tempSelectedRectIndex = ref(null)
const isDrawing = ref(false); // 是否正在绘制矩形框
const undoStack = ref([]); // 撤销操作栈
const currentMode = ref("rectTool"); // 当前模式

const currentClassification = ref(1);
const classifications = ref({
    1: {
        id: 1,
        name: "分类1",
        color: "#FF0000",
    },
});

const createClassificaton = (code) => { };

const switchMode = (mode) => {
    currentMode.value = mode;
};

// 计算 selectedRect
const selectedRect = computed(() => {
    const index = selectedRectIndex.value;

    if (index !== null && index >= 0 && index < boxes.value.length) {

        const stage = stageRef.value.getStage();
        const rect = stage.find("Rect")[index];

        if (rect && rect.x !== undefined && rect.y !== undefined) return rect
    }

    return null;
});

// 加载图像函数
const loadImage = () => {
    const img = new Image();
    img.src = imageUrl;
    img.onload = () => {
        image.value = img; // 确保图像完全加载后赋值
        const stage = stageRef.value.getStage(); // 通过 stageRef 获取 Stage 实例

        // 获取图片的原始宽高
        const imgWidth = img.width;
        const imgHeight = img.height;

        // 计算合适的缩放比例，使图片适应窗口大小
        const stageWidth = stage.width();
        const stageHeight = stage.height();

        // 保持宽高比缩放图片
        let scale = Math.min(stageWidth / imgWidth, stageHeight / imgHeight);

        imageConfig.value = {
            x: (stageWidth - imgWidth * scale) / 2, // 居中显示
            y: (stageHeight - imgHeight * scale) / 2, // 居中显示
            width: imgWidth * scale,
            height: imgHeight * scale,
            image: img, // 将图像绑定到配置中
        };

        stage.batchDraw(); // 强制重新绘制
    };
    img.onerror = (error) => {
        console.error("Image loading error:", error);
    };
};

// 缩放功能
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

// 鼠标按下事件：区分是画框操作还是拖拽画布
const handleMouseDown = (e) => {
    const stage = stageRef.value.getStage();
    const pointerPos = stage.getPointerPosition();

    if (e.target === stage || currentMode.value !== "rectTool") {
        return selectedRectIndex.value = null
    }

    // 获取绝对变换并进行坐标转换
    const transform = stage.getAbsoluteTransform().copy();
    transform.invert();
    const adjustedPointerPos = transform.point(pointerPos);

    // 创建一个新的矩形框，初始化起点坐标
    boxes.value.push({
        classification: currentClassification.value,
        x: adjustedPointerPos.x,
        y: adjustedPointerPos.y,
        width: 0,
        height: 0,
        stroke: classifications.value[currentClassification.value].color,
        strokeWidth: 2,
        draggable: true
    });

    tempSelectedRectIndex.value = boxes.value.length - 1; // 记录当前选中的矩形
    isDrawing.value = true;
};

// 鼠标移动事件：根据当前模式执行操作
const handleMouseMove = () => {
    if (isDrawing.value) {
        const stage = stageRef.value.getStage();
        const pointerPos = stage.getPointerPosition();
        const box = boxes.value[tempSelectedRectIndex.value];

        // 获取绝对变换，用于处理缩放和拖动
        const transform = stage.getAbsoluteTransform().copy();
        transform.invert(); // 获取反向变换

        // 使用反向变换计算真实的坐标
        const adjustedPointerPos = transform.point(pointerPos);

        // 计算矩形的宽高
        box.width = adjustedPointerPos.x - box.x;
        box.height = adjustedPointerPos.y - box.y;

        // 更新矩形
        stage.batchDraw();
    }
};

// 鼠标抬起事件：结束操作
const handleMouseUp = () => {
    if (isDrawing.value) {
        isDrawing.value = false;

        undoStack.value.push([...boxes.value]);
        selectedRectIndex.value = tempSelectedRectIndex.value
    }
};

// 拖动标注框功能
const handleDragMove = (index, event) => {
    console.log(event)
    const box = boxes.value[index];
    box.x = event.target.x();
    box.y = event.target.y();

    const stage = stageRef.value.getStage();
    stage.batchDraw();
};

// 拖动结束事件，保存撤销栈
const handleDragEnd = (index, event) => {
    undoStack.value.push(JSON.parse(JSON.stringify(boxes.value)));
};

// 处理标注框点击事件
const handleRectClick = async (index) => {
    selectedRectIndex.value = index;

    const stage = stageRef.value.getStage();
    const rect = stage.find("Rect")[index];

    console.log(transformerRef.value)

    await nextTick(); // 确保 DOM 更新完成

    if (transformerRef.value && rect) {
        transformerRef.value.nodes([rect]);
        transformerRef.value.getLayer().batchDraw(); // 强制重新绘制
    } else {
        console.error("Transformer or selected rectangle not available.");
    }
};

// 撤销功能
const undo = () => {
    if (undoStack.value.length > 0) {
        undoStack.value.pop(); // 移除当前状态
        const previousState = undoStack.value[undoStack.value.length - 1];

        // 如果有之前的状态，使用新数组替换 boxes.value
        if (previousState) {
            boxes.value = [...previousState]; // 用新的数组引用替换 boxes
        } else {
            boxes.value = []; // 如果撤销到最初状态时，清空所有标注框
        }
        const stage = stageRef.value.getStage(); // 通过 stageRef 获取 Stage 实例
        stage.batchDraw();
    }
};

// 删除选中的矩形框
const deleteRect = () => {
    if (selectedRectIndex.value !== null) {
        boxes.value.splice(selectedRectIndex.value, 1);
        selectedRectIndex.value = null;
        undoStack.value.push([...boxes.value]); // 保存当前状态到撤销栈
    }
};

// 监听键盘事件：快捷键增加和删除标注框
const handleKeydown = (event) => {
    console.log(event);
    if (event.key.toLowerCase() === "n" && mode === "selectBox") {
        // 按下 'n' 键新增一个标注框
        boxes.value.push({
            x: 50,
            y: 50,
            width: 100,
            height: 100,
            stroke: "#FF0000",
            strokeWidth: 2,
        });
        undoStack.value.push([...boxes.value]); // 保存当前状态到撤销栈
    } else if (event.key.toLowerCase() === "delete" || event.key.toLowerCase() === "backspace") {
        deleteRect();
    } else if (
        (event.ctrlKey || event.metaKey) &&
        event.key.toLowerCase() === "z"
    ) {
        undo(); // 按下 Ctrl + Z 撤销
        console.log("ok");
    } else if
        (event.key.toLowerCase() === "a") {
        currentMode.value = "selectBox"
    } else if (event.key.toLowerCase() === "u") {
        currentMode.value = "rectTool"
    }
};

// 在组件挂载时加载图像和绑定键盘事件
onMounted(() => {
    loadImage();
    window.addEventListener("keydown", handleKeydown);
});

// 在组件卸载时移除键盘事件监听器
onUnmounted(() => {
    window.removeEventListener("keydown", handleKeydown);
});
</script>

<template>
    <div>
        <div class="flex justify-between items-center px-4 py-2 border-b">
            <n-button type="info"> 获取新任务 </n-button>
            <n-button strong secondary type="primary"> 保存 </n-button>
        </div>
        <div class="flex justify-between py-4 px-10">
            <div class="flex flex-col gap-4 mt-8">
                <n-button strong secondary :type="currentMode === 'selectBox' ? 'primary' : 'tertiary'"
                    @click="switchMode('selectBox')">
                    选择
                    <template #icon>
                        <n-icon :component="NearMeRound" size="24" />
                    </template>
                </n-button>
                <n-button strong secondary :type="currentMode === 'rectTool' ? 'primary' : 'tertiary'"
                    @click="switchMode('rectTool')">
                    矩形工具
                    <template #icon>
                        <n-icon :component="SquareRound" size="24" />
                    </template>
                </n-button>
            </div>
            <div class="border border-2">
                <v-stage ref="stageRef" :config="stageConfig" @wheel="handleZoom" @mousedown="handleMouseDown"
                    @mousemove="handleMouseMove" @mouseup="handleMouseUp">
                    <v-layer>
                        <v-image :image="image" :config="imageConfig" />
                        <!-- 绘制标注框 -->
                        <v-rect v-for="(rect, index) in boxes" :key="index" :config="rect"
                            @dragmove="handleDragMove(index, $event)" @dragend="handleDragEnd(index, $event)"
                            @click="handleRectClick(index)" />
                        <!-- Transformer 控件，用于调整标注框 -->
                        <v-transformer v-if="selectedRect && selectedRectIndex !== null" ref="transformerRef"
                            :nodes="[selectedRect]" />
                    </v-layer>
                </v-stage>
            </div>
            <div>
                <div class="p-4 mt-8 border border-1 rounded-md">
                    <div>分类</div>

                    <div class="mt-4" v-for="(item, index) in classifications" :key="index">
                        <n-button strong ghost type="error">
                            {{ item.id }} - {{ item.name }}
                            <template #icon>
                                <n-icon :component="Brightness1Round" size="16" :color="item.color" />
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
