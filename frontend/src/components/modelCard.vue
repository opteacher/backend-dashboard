<template>
<el-card v-if="show" class="box-card model-card" style="width: 30vw" :style="`top:${y}px;left:${x}`">
    <div slot="header" class="clearfix">
        <span class="card-name"
            @mousedown="move = true"
            @mouseup="move = false"
            @mousemove="mousemove">{{model.title}}</span>
        <el-link class="float-right" @click="deleteModel">
            <i class="el-icon-close"/>
        </el-link>
    </div>
    <div class="text item">
        tttttt
    </div>
</el-card>
</template>

<script>
export default {
    props: {
        "model": Object
    },
    data() {
        return {
            show: true,
            move: false,
            x: 0,
            y: 0
        }
    },
    async created() {
    },
    methods: {
        deleteModel() {
            this.show = false
        },
        mousemove(me) {
            if (this.move) {
                this.x += me.offsetX;
                this.y += me.offsetY;
            }
        }
    }
}
</script>


<style lang="scss">
.card-name {
    cursor: move;
}

.model-card {
    position: relative;
}
</style>
