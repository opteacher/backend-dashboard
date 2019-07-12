<template>
<dashboard>
    <info-bar @select-api="selectApi"/>
    <div id="pnlFlows" class="w-100 h-100" style="position: absolute"></div>
    <svg id="pnlGraphs" class="w-100" style="position: absolute; z-index: -100; height: 100%" />
    <el-button id="btnApiInfo" type="primary" size="small" icon="el-icon-info" v-if="selApi" @click="showApiInfo">
        {{selApi.name}}
    </el-button>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog :title="`步骤 #${selStep.index}`" :visible.sync="showStepDtlDlg" :modal-append-to-body="false" width="50vw">
        <step-detail ref="step-detail-form" :selStep="selStep"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="chgOperStep">编 辑</el-button>
            <el-button type="primary">确 定</el-button>
        </div>
    </el-dialog>
</dashboard>
</template>

<script>
import _ from "lodash"

import dashboard from "../layouts/dashboard"
import infoBar from "../components/infoBar"
import stepDetail from "../forms/stepDetail"
import backend from "../async/backend"
import apiInfo from "../forms/apiInfo"

export default {
    components: {
        "dashboard": dashboard,
        "info-bar": infoBar,
        "step-detail": stepDetail
    },
    data() {
        return {
            selApi: null,
            selStep: {
                index: -1
            },
            index: 1,
            showStepDtlDlg: false,
            fors: [],
        }
    },
    methods: {
        updatePanel() {
            d3.select("#pnlFlows").html("")
            d3.select("#pnlGraphs").html("")
            this.drawFlowBlock()
            this.drawFlowArrow()
        },
        selectApi(selApi) {
            this.selApi = selApi
            this.selApi.flows = this.selApi.flows.map((flow, idx) => {
                // 做一些处理：只包含一个元素的输出数组全部设为空
                if (flow.outputs.length === 1 && flow.outputs[0] === "") {
                    flow.outputs = []
                }
                if (flow.requires.length === 1 && flow.requires[0] === "") {
                    flow.requires = []
                }
                if (idx === 0) {
                    flow.locVars = _.keys(this.selApi.params)
                } else {
                    flow.locVars = this.selApi.flows[idx - 1].outputs
                }
                flow.index = idx
                return flow
            })
            this.updatePanel()
        },
        drawFlowBlock() {
            let pnlWid = parseInt(document.getElementById("pnlFlows").getBoundingClientRect().width)
            let flowLoc = 50
            let flowX = (pnlWid>>1) - 250
            let card = d3.select("#pnlFlows")
                .selectAll("div")
                .data(this.selApi.flows)
                .join("div")
                .attr("class", "card")
                .attr("name", (flow, idx) => `flow_${idx}`)
                .style("position", "absolute")
                .style("left", flow => `${flow.x = flowX}px`)
                .style("top", (flow, idx) => `${flow.y = (idx === 0 ? flowLoc : flowLoc += 200)}px`)
                .style("width", "500px")
                .style("margin-bottom", (flow, idx) => `${idx === this.selApi.flows.length - 1 ? 50 : 0}px`)
                .each(flow => {
                    if (!flow.special) {
                        return
                    }
                    // 收集步骤的特殊标识
                    switch (flow.special) {
                        case 1:// 循环开始标识
                            this.fors.push({
                                begin: flow
                            })
                            break
                        case 2:// 循环结束标识
                            // 找出没有结束标识的begin，并用离这个结束块最近的作为end
                            let noEnds = this.fors.filter(f => !f.end)
                            let clsEnd = noEnds[0]
                            for (let fblk of noEnds) {
                                if (fblk.begin.index < clsEnd.begin.index) {
                                    clsEnd = fblk
                                }
                            }
                            clsEnd.end = flow
                            break
                    }
                })
                .append("div")
                .attr("class", "row")
            // 填充步骤的inputs
            card.append("div")
                .attr("class", "col pr-0")
                .append("div")
                .attr("class", "list-group list-group-flush h-100")
                .selectAll("a")
                .data((flow, idx) => _.toPairs(flow.inputs).map(kv => {
                    return {
                        pholder: kv[0],
                        content: kv[1],
                        findex: idx
                    }
                }))
                .join("a")
                .attr("class", "list-group-item list-group-item-primary list-group-item-action api-params")
                .attr("href", "#")
                .text(input => input.pholder)
                .append("i")
                .attr("class", "el-icon-arrow-right")
            // 填充步骤的描述
            card.append("div")
                .attr("class", "col-6 card-body text-center desc-panel")
                .text(flow => flow.desc)
                .on("click", flow => {
                    this.selStep = flow
                    this.showStepDtlDlg = true
                })
            // 填充步骤的outputs
            card.append("div")
                .attr("class", "col pl-0")
                .append("div")
                .attr("class", "list-group list-group-flush h-100")
                .selectAll("a")
                .data(flow => flow.outputs)
                .join("a")
                .attr("class", "list-group-item list-group-item-success list-group-item-action api-params text-right")
                .attr("href", "#")
                .text(output => output)
                .append("i")
                .attr("class", "el-icon-arrow-right")
            // 绘制局部变量列表
            card.append("div")
                .attr("class", "list-group")
                .style("position", "absolute")
                .style("left", flow => `-${flow.x + 2}px`)
                .style("top", 0)
                .selectAll("a")
                .data(flow => flow.locVars)
                .join("a")
                .attr("class", "list-group-item list-group-item-action local-vars rl-0")
                .text(v => v)
                .each(function(v) {
                    let rect = this.getBoundingClientRect()
                    let x1 = rect.width
                    let y1 = rect.y + (rect.height>>1)
                    // @_@
                })
        },
        drawFlowArrow() {
            let self = this
            d3.select("#pnlGraphs")
                .style("height", `${document.getElementById("pnlFlows").scrollHeight}px`)
                .selectAll("g")
                .data(this.selApi.flows)
                .join("line")
                .attr("stroke-width", 2)
                .attr("stroke", "black")
                .each(function(flow, idx) {
                    if (idx === self.selApi.flows.length - 1) {
                        return
                    }
                    let rect = document.getElementsByName(`flow_${idx}`)[0].getBoundingClientRect()
                    let next = self.selApi.flows[idx + 1]
                    let x1 = flow.x + (rect.width>>1)
                    let y1 = flow.y + rect.height
                    let x2 = next.x + (rect.width>>1)
                    let y2 = next.y
                    d3.select(this)
                        .attr("name", `line_${idx}_${idx + 1}`)
                        .attr("x1", x1)
                        .attr("y1", y1)
                        .attr("x2", x2)
                        .attr("y2", y2)
                    // 画箭头
                    d3.select("#pnlGraphs")
                        .append("polyline")
                        .attr("fill", "black")
                        .attr("stroke", "blue")
                        .attr("stroke-width", 2)
                        .attr("points", [
                            `${x2 - 5},${next.y - 10}`,
                            `${x2},${next.y}`,
                            `${x2 + 5},${next.y - 10}`
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
                    let x = ((x2 - x1)>>1) + x1
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
            // 折线的步进
            let stepFor = 20
            // 绘制循环的折线箭头
            d3.select("#pnlGraphs")
                .selectAll("g")
                .data(this.fors)
                .join("polyline")
                .attr("fill", "none")
                .attr("stroke-width", 2)
                .attr("stroke", "black")
                .each(function(forBlk) {
                    let forBeg = forBlk.begin
                    let begHftHgt = parseInt(document.getElementsByName(`flow_${forBeg.index}`)[0].getBoundingClientRect().height)>>1
                    let forEnd = forBlk.end
                    let endHftHgt = parseInt(document.getElementsByName(`flow_${forEnd.index}`)[0].getBoundingClientRect().height)>>1
                    let width = 500
                    let points = [
                        `${forEnd.x + width},${forEnd.y + endHftHgt}`,
                        `${forEnd.x + width + stepFor},${forEnd.y + endHftHgt}`,
                        `${forBeg.x + width + stepFor},${forBeg.y + begHftHgt}`,
                        `${forBeg.x + width},${forBeg.y + begHftHgt}`,
                    ]
                    d3.select(this).attr("points", points.join(" "))
                    // 绘制箭头
                    let ex = forBeg.x + width
                    let ey = forBeg.y + begHftHgt
                    d3.select("#pnlGraphs")
                        .append("polyline")
                        .attr("fill", "black")
                        .attr("stroke", "blue")
                        .attr("stroke-width", 2)
                        .attr("points", [
                            `${ex + 10},${ey + 5}`,
                            `${ex},${ey}`,
                            `${ex + 10},${ey - 5}`,
                        ].join(" "))
                })
        },
        drawCurve(x1, y1, x2, y2, dir1, dir2) {
            let data = [{
                x: x1, y: y1
            }, {
                x: x2, y: y2
            }]
            let func = d3.svg.line()
                .x(d => d.x).y(d => d.y)
                .interpolate("basis")
            d3.select("#pnlGraphs")
                .append("path")
                .attr("d", func(data))
                .attr("stroke", "blue")
                .attr("stroke-width", 1)
                .attr("fill", "none")
        },
        chgOperStep() {
            this.$refs["step-detail-form"].mode = (
                this.$refs["step-detail-form"].mode === "display" ? "editing" : "display"
            )
        },
        showApiInfo() {
            this.index++
            this.$msgbox({
                title: "接口信息",
                message: this.$createElement(apiInfo, {
                    props: {
                        api: this.selApi,
                    },
                    key: this.index
                }),
                showConfirmButton: false
            }).catch(err => {})
        },
    }
}
</script>

<style lang="scss">
.api-params, .local-vars {
    font-size: 0.2rem;
    padding: .5vh .5vw;
}
.desc-panel:hover {
    cursor: pointer;
    background-color: #f8f9fa;
}
#btnApiInfo {
    position: fixed;
    right: 0;
    bottom: 0;
    border-radius: 4px 0 0 0;
}
</style>
