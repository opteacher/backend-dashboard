<template>
<g v-if="istFlowBtn.next">
    <line :name="`flowLink${istFlowBtn.nsuffix}`" stroke-width="2" stroke="black"/>
    <polygon :name="`flowArrow${istFlowBtn.nsuffix}`" fill="black" stroke="blue" stroke-width="2"/>
    <line :name="`flowSplit${istFlowBtn.nsuffix}`" x1="0" stroke="black" stroke-dasharray="5,5"/>
</g>
</template>

<script>
export default {
    props: {
        "istFlowBtn": Object
    },
    mounted() {
        d3.select("#pnlGraphs")
            .style("height", `${document.getElementById("pnlFlows").scrollHeight}px`)

        let prev = this.istFlowBtn.prev
        // 如果是结尾标识，则不继续显示下面的内容
        if (prev.symbol === 3/* SpcSymbol.END */) {
            d3.select(`button[name="istFlowBtn${this.istFlowBtn.nsuffix}"]`)
                .style("display", "none")
            return
        }
        let prevName = `flow_${prev.index}`
        let prevRect = document.getElementsByName(prevName)[0].getBoundingClientRect()
        let x1 = prev.x + (prevRect.width>>1)
        let y1 = prev.y + prevRect.height

        let next = this.istFlowBtn.next
        let nextRect = prevRect //@_@
        let x2 = next ? (next.x + (nextRect.width>>1)) : x1
        let y2 = next ? next.y : 200

        let x = ((x2 - x1)>>1) + x1
        let y = next ? (((y2 - y1)>>1) + y1) : y2

        // 画连线
        d3.select(`line[name="flowLink${this.istFlowBtn.nsuffix}"]`)
            .attr("x1", x1).attr("y1", y1)
            .attr("x2", x2).attr("y2", y2)
            
        // 画箭头
        d3.select(`polygon[name="flowArrow${this.istFlowBtn.nsuffix}"]`)
            .attr("points", [
                `${x2 - 5},${next.y - 10}`,
                `${x2},${next.y}`,
                `${x2 + 5},${next.y - 10}`
            ].join(" "))
        
        // 画分隔线
        d3.select(`line[name="flowSplit${this.istFlowBtn.nsuffix}"]`)
            .attr("y1", y)
            .attr("x2", document.getElementById("pnlFlows").getBoundingClientRect().width)
            .attr("y2", y)

        // 调整插入步骤按钮
        d3.select(`button[name="istFlowBtn${this.istFlowBtn.nsuffix}"]`)
            .style("left", `${x - 20}px`)
            .style("top", `${y - 20}px`)
    }
}
</script>

