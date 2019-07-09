<template>
<dashboard>
    <info-bar @sel-interface="selInterface"/>
    <div id="pnlFLows" class="w-100 h-100" style="position: absolute"></div>
    <svg id="pnlGraphs" class="w-100" style="position: absolute; z-index: -100; height: 100%" />
</dashboard>
</template>

<script>
import dashboard from "../layouts/dashboard"
import infoBar from "../components/infoBar"

export default {
    components: {
        "dashboard": dashboard,
        "info-bar": infoBar
    },
    data() {
        return {
            selItf: null
        }
    },
    methods: {
        selInterface(selItf) {
            this.selItf = selItf
            this.drawFlowBlock()
            this.drawFlowArrow()
        },
        drawFlowBlock() {
            let pnlWid = parseInt(document.getElementById("pnlFLows").getBoundingClientRect().width)
            let flowLoc = 50
            let flowX = (pnlWid>>1) - 150
            let card = d3.select("#pnlFLows")
                .html("")
                .selectAll("div")
                .data(this.selItf.flows)
                .join("div")
                .attr("class", "card")
                .attr("name", (flow, idx) => `flow_${idx}`)
                .style("position", "absolute")
                .style("cursor", "pointer")
                .style("left", flow => `${flow.x = flowX}px`)
                .style("top", (flow, idx) => {
                    flow.y = (idx === 0 ? flowLoc : flowLoc += 200)
                    return `${flow.y}px`
                })
                .style("width", "300px")
                .style("margin-bottom", (flow, idx) => `${idx === this.selItf.flows.length - 1 ? 50 : 0}px`)
                .append("div")
                .attr("class", "card-body")
                .text(flow => flow.desc)
            card.append("a")
                .attr("href", "#")
                .style("position", "absolute")
                .style("left", 0)
                .style("top", 0)
                .attr("class", "badge badge-primary")
                .text("abcd")
                .append("i")
                .attr("class", "el-icon-arrow-right")
                .each(function() {
                    console.log(this.getBoundingClientRect())
                })
        },
        drawFlowArrow() {
            let self = this
            let pnlHgt = document.getElementById("pnlFLows").scrollHeight
            d3.select("#pnlGraphs")
                .html("")
                .style("height", `${pnlHgt}px`)
                .selectAll("g")
                .data(this.selItf.flows)
                .join("line")
                .attr("stroke-width", 2)
                .attr("stroke", "black")
                .each(function(flow, idx) {
                    if (idx === self.selItf.flows.length - 1) {
                        return
                    }
                    let rect = document.getElementsByName(`flow_${idx}`)[0].getBoundingClientRect()
                    let next = self.selItf.flows[idx + 1]
                    let x = flow.x + (rect.width>>1)
                    let y1 = flow.y + rect.height
                    let y2 = next.y
                    d3.select(this)
                        .attr("name", `line_${idx}_${idx + 1}`)
                        .attr("x1", x)
                        .attr("y1", y1)
                        .attr("x2", x)
                        .attr("y2", y2)
                    d3.select("#pnlGraphs")
                        .append("polyline")
                        .attr("fill", "black")
                        .attr("stroke", "blue")
                        .attr("stroke-width", 2)
                        .attr("points", [
                            `${x - 5},${next.y - 10}`,
                            `${x},${next.y}`,
                            `${x + 5},${next.y - 10}`
                        ].join(" "))
                    // 按钮宽高40px
                    d3.select("#pnlFLows")
                        .append("button")
                        .attr("class", "btn btn-success rounded-circle")
                        .attr("type", "button")
                        .style("position", "absolute")
                        .style("left", `${x - 20}px`)
                        .style("top", `${((y2 - y1)>>1) + y1 - 20}px`)
                        .append("i")
                        .attr("class", "el-icon-plus")
                })
        }
    }
}
</script>
