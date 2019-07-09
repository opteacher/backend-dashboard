<template>
<dashboard>
    <info-bar @sel-interface="selInterface"/>
    <div id="pnlFlows" class="w-100 h-100" style="position: absolute"></div>
    <svg id="pnlGraphs" class="w-100" style="position: absolute; z-index: -100; height: 100%" />
</dashboard>
</template>

<script>
import _ from "lodash"

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
            let pnlWid = parseInt(document.getElementById("pnlFlows").getBoundingClientRect().width)
            let flowLoc = 50
            let flowX = (pnlWid>>1) - 250
            let card = d3.select("#pnlFlows")
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
                .style("width", "500px")
                .style("margin-bottom", (flow, idx) => `${idx === this.selItf.flows.length - 1 ? 50 : 0}px`)
                .append("div")
                .attr("class", "row")
            card.append("div")
                .attr("class", "col pr-0")
                .append("div")
                .attr("class", "list-group list-group-flush")
                .selectAll("a")
                .data(flow => _.toPairs(flow.inputs))
                .join("a")
                .attr("class", "list-group-item list-group-item-action api-params")
                .attr("href", "#")
                .text(input => input[0])
                .append("i")
                .attr("class", "el-icon-arrow-right")
            card.append("div")
                .attr("class", "col-6 card-body text-center")
                .text(flow => flow.desc)
            card.append("div")
                .attr("class", "col pl-0")
                .append("div")
                .attr("class", "list-group list-group-flush")
                .selectAll("a")
                .data(flow => flow.outputs)
                .join("a")
                .attr("class", "list-group-item list-group-item-action api-params text-right")
                .attr("href", "#")
                .text(output => output)
                .append("i")
                .attr("class", "el-icon-arrow-right")
            d3.select("#pnlFlows")
                .append("div")
                .attr("class", "list-group")
                .style("position", "absolute")
                .style("left", 0)
                .style("top", "50px")
                .selectAll("a")
                .data(_.toPairs(this.selItf.params))
                .join("a")
                .attr("class", "list-group-item list-group-item-action local-vars")
                .text(param => param[0])
        },
        drawFlowArrow() {
            let self = this
            let pnlHgt = document.getElementById("pnlFlows").scrollHeight
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
                    // 画箭头
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
                    // 画步骤分隔线
                    let y = ((y2 - y1)>>1) + y1
                    d3.select("#pnlGraphs")
                        .append("line")
                        .attr("x1", 0)
                        .attr("y1", y)
                        .attr("x2", document.getElementById("pnlGraphs").getBoundingClientRect().width)
                        .attr("y2", y)
                        .attr("stroke", "black")
                        .attr("stroke-dasharray","5,5")
                    // 按钮宽高40px
                    d3.select("#pnlFlows")
                        .append("button")
                        .attr("class", "btn btn-success rounded-circle")
                        .attr("type", "button")
                        .style("position", "absolute")
                        .style("left", `${x - 20}px`)
                        .style("top", `${y - 20}px`)
                        .append("i")
                        .attr("class", "el-icon-plus")
                })
        }
    }
}
</script>

<style lang="scss">
.api-params, .local-vars {
    font-size: 0.2rem;
    padding: .5vh .5vw;
}
</style>
