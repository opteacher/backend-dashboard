<template>
<line :x1="x1" :y1="y1" :x2="x2" :y2="y2" stroke="#000" stroke-width="2"/>
</template>

<script>
export default {
    props: {
        "relation": Object,
        "begModel": Object,
        "endModel": Object,
    },
    data() { return {
        x1: 0,
        y1: 0,
        x2: 0,
        y2: 0
    }},
    created() {
        console.log(this.begModel)
        console.log(this.endModel)
    }
}
</script>

