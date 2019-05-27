<template>
<canvas id="pnlModels"></canvas>
</template>

<script>
export default {
    data() {
        return {
            models: []
        }
    },
    created() {
        let canvas = document.getElementById("pnlModels")
        this.stage = new JTopo.Stage(canvas)
        this.scene = new JTopo.Scene(this.stage)

        this.node = new JTopo.Node("Hello")
        this.node.setLocation(200, 300)
        this.scene.add(this.node)
    }
}
</script>
