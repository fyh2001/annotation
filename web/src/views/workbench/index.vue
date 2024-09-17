<template>
    <v-stage ref="stageRef" :config="stageConfig" @wheel="handleZoom" @mousedown="handleMouseDown"
        @mousemove="handleMouseMove" @mouseup="handleMouseUp">
        <v-layer>
            <v-image :image="image" />
            <!-- 绘制标注框 -->
            <v-rect v-for="(rect, index) in boxes" :key="index" :config="rect" draggable
                @dragmove="handleDragMove(index, $event)" @click="handleRectClick(index)" />
            <!-- Transformer 控件，用于调整标注框 -->
            <v-transformer v-if="selectedRectIndex !== null" ref="transformer" :nodes="[selectedRect]" />
        </v-layer>
    </v-stage>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import imageUrl from "../../../public/1.jpg"

// 初始化 Konva 画布配置
const stageConfig = ref({
    width: window.innerWidth,
    height: window.innerHeight,
    draggable: true,  // 初始状态下允许画布拖动
    scaleX: 1,
    scaleY: 1,
});

const stageRef = ref(null);
const transformer = ref(null);

// 存储图像、标注框和缩放比例
const image = ref(null);  // 存放加载的图像
const boxes = ref([]);    // 存放标注框
const scale = ref(1);     // 缩放比例
const selectedRectIndex = ref(null);  // 当前选中的矩形框索引
const isDrawing = ref(false);  // 是否正在绘制矩形框
const isDraggingCanvas = ref(false);  // 是否正在拖拽画布
const newRect = ref(null);    // 当前正在绘制的新矩形框
const undoStack = ref([]);    // 撤销操作栈

// 加载图像函数
const loadImage = () => {
    const img = new Image();
    img.src = imageUrl;
    img.onload = () => {
        image.value = img;  // 确保图像完全加载后赋值
        const stage = stageRef.value.getStage();  // 通过 stageRef 获取 Stage 实例
        stage.batchDraw();  // 强制重新绘制
    };
    img.onerror = (error) => {
        console.error('Image loading error:', error);
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

    if (e.target === stage) {
        // 如果点击的是空白区域，启用画布拖拽模式
        isDrawing.value = false;
        isDraggingCanvas.value = true;
        stageConfig.value.draggable = true;
    } else {
        // 点击非空白区域时，启用画框模式
        isDrawing.value = true;
        isDraggingCanvas.value = false;
        stageConfig.value.draggable = false;  // 禁用画布拖拽
        newRect.value = {
            // x: pointerPos.x,
            // y: pointerPos.y,
            x: stage.getRelativePointerPosition().x,
            y: stage.getRelativePointerPosition().y,
            width: 0,
            height: 0,
            stroke: 'red',
            strokeWidth: 2,
        };
        boxes.value.push(newRect.value);
        selectedRectIndex.value = boxes.value.length - 1;

    }
};

// 鼠标移动事件：根据当前模式执行操作
const handleMouseMove = () => {
    if (isDrawing.value) {
        const stage = stageRef.value.getStage();
        const pointerPos = stage.getPointerPosition();
        const box = boxes.value[selectedRectIndex.value];

        // 计算当前缩放下的坐标
        const scale = stage.scaleX();
        box.width = (pointerPos.x - box.x) / scale;
        box.height = (pointerPos.y - box.y) / scale;

        // 重新设置矩形的大小
        stage.batchDraw();
    }
};

// 鼠标抬起事件：结束操作
const handleMouseUp = () => {
    if (isDrawing.value) {
        isDrawing.value = false;
        undoStack.value.push([...boxes.value]);  // 保存当前状态到撤销栈
    }
    if (isDraggingCanvas.value) {
        isDraggingCanvas.value = false;
    }
    stageRef.value.getStage().config.draggable = true;  // 重新启用画布拖拽
};

// 拖动标注框功能
const handleDragMove = (index, event) => {
    const box = boxes.value[index];
    box.x = event.target.x();
    box.y = event.target.y();
    undoStack.value.push([...boxes.value]);  // 保存当前状态到撤销栈
};

// 处理标注框点击事件
const handleRectClick = (index) => {
    selectedRectIndex.value = index;
    transformer.value.nodes([stageRef.value.getStage().find('Rect')[index]]);
};

// 撤销功能
const undo = () => {
    if (undoStack.value.length > 0) {
        boxes.value = undoStack.value.pop();
    }
};

// 删除选中的矩形框
const deleteRect = () => {
    if (selectedRectIndex.value !== null) {
        boxes.value.splice(selectedRectIndex.value, 1);
        selectedRectIndex.value = null;
        undoStack.value.push([...boxes.value]);  // 保存当前状态到撤销栈
    }
};

// 监听键盘事件：快捷键增加和删除标注框
const handleKeydown = (event) => {
    if (event.key === 'n') {
        // 按下 'n' 键新增一个标注框
        boxes.value.push({
            x: 50,
            y: 50,
            width: 100,
            height: 100,
            stroke: 'red',
            strokeWidth: 2,
        });
        undoStack.value.push([...boxes.value]);  // 保存当前状态到撤销栈
    } else if (event.key === 'Delete') {
        deleteRect();
    } else if (event.ctrlKey && event.key === 'z') {
        undo();  // 按下 Ctrl + Z 撤销
    }
};

// 在组件挂载时加载图像和绑定键盘事件
onMounted(() => {
    loadImage();
    window.addEventListener('keydown', handleKeydown);
});

// 在组件卸载时移除键盘事件监听器
onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown);
});
</script>

<style>
body {
    margin: 0;
    overflow: hidden;
}
</style>