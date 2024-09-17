<template>
    <div>
        <v-stage ref="stage" :config="stageConfig" @wheel="handleZoom" @mousedown="handleStageClick">
            <v-layer ref="layer">
                <!-- 矩形标注框 -->
                <v-rect v-for="(box, index) in boxes" :key="index" :config="box" @click="handleRectClick(index)"
                    :draggable="true" @dragmove="handleDragMove(index)" />

                <!-- Transformer 控件，用于调整标注框 -->
                <v-transformer v-if="selectedRectIndex !== null" ref="transformer" />
            </v-layer>
        </v-stage>
    </div>
</template>

<script>
import { Stage, Layer, Rect, Transformer } from 'vue-konva';

export default {
    components: {
        'v-stage': Stage,
        'v-layer': Layer,
        'v-rect': Rect,
        'v-transformer': Transformer,
    },
    data() {
        return {
            stageConfig: {
                width: window.innerWidth,
                height: window.innerHeight,
                draggable: true,
                scaleX: 1,
                scaleY: 1,
            },
            boxes: [], // 存储所有的矩形标注框
            selectedRectIndex: null, // 当前选中的矩形索引
        };
    },
    methods: {
        // 处理缩放
        handleZoom(e) {
            e.evt.preventDefault();
            const scaleBy = 1.05;
            const stage = this.$refs.stage.getStage();
            const oldScale = stage.scaleX();
            const pointer = stage.getPointerPosition();

            const direction = e.evt.deltaY > 0 ? 1 : -1;
            const newScale = direction > 0 ? oldScale * scaleBy : oldScale / scaleBy;

            stage.scale({ x: newScale, y: newScale });
            const newPos = {
                x: pointer.x - (pointer.x - stage.x()) * (newScale / oldScale),
                y: pointer.y - (pointer.y - stage.y()) * (newScale / oldScale),
            };
            stage.position(newPos);
            stage.batchDraw();
        },

        // 处理矩形点击，选择并附加 Transformer
        handleRectClick(index) {
            this.selectedRectIndex = index;
            this.updateTransformer();
        },

        // 更新 Transformer 的附加对象
        updateTransformer() {
            const stage = this.$refs.stage.getStage();
            const transformer = this.$refs.transformer.getNode();
            const selectedRect = stage.find('Rect')[this.selectedRectIndex];

            transformer.nodes([selectedRect]);
            transformer.getLayer().batchDraw();
        },

        // 处理舞台点击（用于取消选中矩形）
        handleStageClick(e) {
            if (e.target === e.target.getStage()) {
                this.selectedRectIndex = null;
                this.$refs.transformer.getNode().detach();
                this.$refs.layer.getNode().batchDraw();
            }
        },

        // 处理矩形的拖拽移动
        handleDragMove(index) {
            const stage = this.$refs.stage.getStage();
            const newPos = stage.find('Rect')[index].position();

            this.boxes[index].x = newPos.x;
            this.boxes[index].y = newPos.y;
        },
    },
};
</script>

<style scoped>
body {
    margin: 0;
    padding: 0;
    overflow: hidden;
}
</style>