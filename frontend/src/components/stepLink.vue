<template>
<g>
    <line :name="`stepLink${istStepBtn.nsuffix}`" stroke-width="2" stroke="black"/>
    <polygon :name="`stepArrow${istStepBtn.nsuffix}`" fill="black" stroke="blue" stroke-width="2"/>
    <line :name="`stepSplit${istStepBtn.nsuffix}`" x1="0" stroke="black" stroke-dasharray="5,5"/>
    <rect :name="`stepTest${istStepBtn.nsuffix}`" fill="green"/>
</g>
</template>

<script>
export default {
    props: {
        "istStepBtn": Object
    },
    mounted() {
        let pnlFlows = document.getElementById("pnlFlows")
        let pnlRect = pnlFlows.getBoundingClientRect()
        if (!this.istStepBtn.prev && !this.istStepBtn.next) {
            // 没有步骤，只有一个按钮供添加
            d3.select(`button[name="istStepBtn${this.istStepBtn.nsuffix}"]`)
                .style("left", `${(pnlRect.width>>1) - 20}px`)
                .style("top", "50px")
            return
        }
        
        d3.select("#pnlGraphs")
            .style("height", `${pnlFlows.scrollHeight}px`)

        let prev = this.istStepBtn.prev
        // 如果是结尾标识，则不继续显示下面的内容
        if (prev.symbol & 4 /* SpcSymbol_END */ !== 0) {
            d3.select(`button[name="istStepBtn${this.istStepBtn.nsuffix}"]`).remove()
            return
        }
        let prevName = `step_${prev.index}`
        let prevRect = document.getElementsByName(prevName)[0].getBoundingClientRect()
        let px = prevRect.x - pnlRect.x
        let py = prevRect.y - pnlRect.y
        let x1 = px + (prevRect.width>>1)
        let y1 = py + prevRect.height

        let next = this.istStepBtn.next
        let nx, ny = 0
        let x2, y2 = 0
        let x, y = 0
        if (next) {
            let nextName = `step_${next.index}`
            let nextRect = document.getElementsByName(nextName)[0].getBoundingClientRect()
            nx = nextRect.x - pnlRect.x
            ny = nextRect.y - pnlRect.y
            x2 = nx + (nextRect.width>>1)
            y2 = ny

            x = ((x2 - x1)>>1) + x1
            y = ((y2 - y1)>>1) + y1
        } else {
            x2 = x1
            y2 = y1 + 80 // 这里的80是最后块和添加按钮之间连线的长度

            x = x2
            y = y2
            y2 -= 20
        }

        // 画连线
        d3.select(`line[name="stepLink${this.istStepBtn.nsuffix}"]`)
            .attr("x1", x1).attr("y1", y1)
            .attr("x2", x2).attr("y2", y2)
            
        if (next) {
            // 画箭头
            d3.select(`polygon[name="stepArrow${this.istStepBtn.nsuffix}"]`)
                .attr("points", [
                    `${x2 - 5},${ny - 10}`,
                    `${x2},${ny}`,
                    `${x2 + 5},${ny - 10}`
                ].join(" "))
            
            // 画分隔线
            d3.select(`line[name="stepSplit${this.istStepBtn.nsuffix}"]`)
                .attr("y1", y)
                .attr("x2", document.getElementById("pnlFlows").getBoundingClientRect().width)
                .attr("y2", y)
        }
        
        // 调整插入步骤按钮
        d3.select(`button[name="istStepBtn${this.istStepBtn.nsuffix}"]`)
            .style("left", `${x - 20}px`)
            .style("top", `${y - 20}px`)
    }
}
</script>

