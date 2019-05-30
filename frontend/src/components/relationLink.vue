<template>
<svg class="w-100 h-100">
    <line :x1="x1" :y1="y1" :x2="x2" :y2="y2" stroke="#000" stroke-width="2"/>
</svg>
</template>

<script>
export default {
    props: {
        "relation": Object,
        "model1": Object,
        "model2": Object,
    },
    data() { return {
        x1: 0,
        y1: 0,
        x2: 0,
        y2: 0
    }},
    created() {
        this.relation.onModelChanged = () => {
            let pos1 = this.centerPos(this.model1)
            this.x1 = pos1.x
            this.y1 = pos1.y
            let pos2 = this.centerPos(this.model2)
            this.x2 = pos2.x
            this.y2 = pos2.y
        }
        this.relation.onModelChanged()
    },
    methods: {
        centerPos(model) {
            return {
                x: model.x + (model.width>>1),
                y: model.y + (model.height>>1)
            }
        }
    }
}
</script>

