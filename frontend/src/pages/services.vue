<template>
<dashboard>
    <info-bar @select-api="selectApi" @add-api="addApi"/>
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
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog :title="`添加步骤 #${flowInfo.flowIdx}`" :visible.sync="showAddFlowDlg" :modal-append-to-body="false" width="50vw">
        <edit-flow ref="add-flow-form" :flowInfo="flowInfo"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddFlowDlg = false">取 消</el-button>
            <el-button type="primary" @click="addFlow">确 定</el-button>
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
import editFlow from "../forms/editFlow"
import stepBkd from "../async/step"
import apisBkd from "../async/api"

export default {
    components: {
        "dashboard": dashboard,
        "info-bar": infoBar,
        "step-detail": stepDetail,
        "edit-flow": editFlow
    },
    data() {
        return {
            selApi: null,
            selStep: {
                index: -1
            },
            index: 1,
            showStepDtlDlg: false,
            showAddFlowDlg: false,
            fors: [],
            istFlowBtns: [],
            flowInfo: {
                flowIdx: -1
            }
        }
    },
    methods: {
        async refresh() {
            let res = await apisBkd.qryByName(this.selApi.name)
            if (typeof res === "string") {
                this.$message.error(`查询名为${this.selApi.name}的接口时发生错误：${res}`)
            } else {
                this.selectApi(res.data.data)
            }
        },
        updatePanel() {
            d3.select("#pnlFlows").html("")
            d3.select("#pnlGraphs").html("")
            this.drawFlowBlock()
            this.drawFlowArrow()
        },
        selectApi(selApi) {
            this.selApi = selApi
            this.istFlowBtns = [{
                apiName: selApi.name,
                flowIdx: 0,
                stepId: -1,
                locVars: this.selApi.params ? _.keys(this.selApi.params) : []
            }]
            this.selApi.flows = this.selApi.flows ? this.selApi.flows.map((flow, idx) => {
                // 做一些处理：只包含一个元素的输出数组全部设为空
                if (!flow.outputs || (flow.outputs.length === 1 && flow.outputs[0] === "")) {
                    flow.outputs = []
                }
                if (!flow.requires || (flow.requires.length === 1 && flow.requires[0] === "")) {
                    flow.requires = []
                }
                flow.index = idx

                if (idx !== 0) {
                    let prevLocVars = this.istFlowBtns[this.istFlowBtns.length - 1].locVars
                    let prevOutputs = this.selApi.flows[idx - 1].outputs
                    this.istFlowBtns.push({
                        apiName: selApi.name,
                        flowIdx: idx + 1,
                        stepId: flow.id,
                        locVars: prevLocVars.concat(prevOutputs)
                    })
                }
                return flow
            }) : []
            this.updatePanel()
        },
        drawFlowBlock() {
            if (!this.selApi.flows || this.selApi.flows.length === 0) {
                // 如果没有处理流程，绘制一个添加步骤的按钮之后直接返回
                d3.select("#pnlFlows")
                    .append("div")
                    .attr("class", "h-100 w-100")
                    .attr("style", "display: flex")
                    .append("button")
                    .attr("style", "align-self: center; margin: 0 auto")
                    .attr("class", "btn btn-success rounded-circle")
                    .attr("type", "button")
                    .on("click", () => {
                        this.flowInfo = this.istFlowBtns[0]
                        this.showAddFlowDlg = true
                    })
                    .append("i")
                    .attr("class", "el-icon-plus")
                return
            }
            let pnlWid = parseInt(document.getElementById("pnlFlows").getBoundingClientRect().width)
            let flowLoc = 50
            let flowX = (pnlWid>>1) - 300
            let disBetwFlow = 300
            let card = d3.select("#pnlFlows")
                .selectAll("div")
                .data(this.selApi.flows)
                .join("div")
                .attr("class", "card")
                .attr("name", (flow, idx) => `flow_${idx}`)
                .style("position", "absolute")
                .style("left", flow => `${flow.x = flowX}px`)
                .style("top", (flow, idx) => `${flow.y = (idx === 0 ? flowLoc : flowLoc += disBetwFlow)}px`)
                .style("width", "600px")
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
            card.append("div")
                .attr("class", "card-header")
                .text((flow, index) => `#${index}`)
                .append("span")
                .attr("class", "float-right")
                .text(flow => flow.operKey)
            let cardBody = card.append("div")
                .attr("class", "row") 
            // 填充步骤的inputs
            cardBody.append("div")
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
            cardBody.append("div")
                .attr("class", "col-6 card-body text-center desc-panel")
                .text(flow => flow.desc)
                .on("click", flow => {
                    this.selStep = flow
                    this.showStepDtlDlg = true
                })
            // 填充步骤的outputs
            cardBody.append("div")
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
        },
        drawFlowArrow() {
            if (!this.selApi.flows || this.selApi.flows.length === 0) {
                // 如果没有处理流程，直接返回
                return
            }
            let self = this
            d3.select("#pnlGraphs")
                .style("height", `${document.getElementById("pnlFlows").scrollHeight}px`)
                .selectAll("g")
                .data(this.selApi.flows)
                .join("line")
                .attr("stroke-width", 2)
                .attr("stroke", "black")
                .each(function(flow, idx) {
                    // 如果存在return特殊标识，则返回
                    if (self.selApi.flows.length !== 1 && idx === self.selApi.flows.length - 1) {
                        return
                    }
                    let rect = document.getElementsByName(`flow_${idx}`)[0].getBoundingClientRect()
                    let next = self.selApi.flows.length === 1 ? null : self.selApi.flows[idx + 1]
                    let x1 = flow.x + (rect.width>>1)
                    let y1 = flow.y + rect.height
                    let x2 = next ? (next.x + (rect.width>>1)) : x1
                    let y2 = next ? next.y : 200
                    d3.select(this)
                        .attr("name", `line_${idx}_${idx + 1}`)
                        .attr("x1", x1)
                        .attr("y1", y1)
                        .attr("x2", x2)
                        .attr("y2", y2)
                    let x = ((x2 - x1)>>1) + x1
                    let y = next ? (((y2 - y1)>>1) + y1) : y2
                    if (next) {
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
                        d3.select("#pnlGraphs")
                            .append("line")
                            .attr("x1", 0)
                            .attr("y1", y)
                            .attr("x2", document.getElementById("pnlGraphs").getBoundingClientRect().width)
                            .attr("y2", y)
                            .attr("stroke", "black")
                            .attr("stroke-dasharray","5,5")
                    }
                    // 绘制添加步骤按钮，按钮宽高40px
                    d3.select("#pnlFlows")
                        .append("button")
                        .attr("class", "btn btn-success rounded-circle")
                        .attr("type", "button")
                        .style("position", "absolute")
                        .style("left", `${x - 20}px`)
                        .style("top", `${y - 20}px`)
                        .on("click", () => {
                            self.flowInfo = self.istFlowBtns.find(ele => ele.flowIdx === idx)
                            self.showAddFlowDlg = true
                        })
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
        addApi(newApi) {
            console.log(newApi)
        },
        async addFlow() {
            let form = this.$refs["add-flow-form"]
            let res = await stepBkd.addFlow({
                index: form.flowInfo.flowIdx,
                operStep: Object.assign(form.flow, {
                    apiName: form.flowInfo.apiName
                })
            })
            if (typeof res === "string") {
                this.$message.error(`插入流程发生错误：${res}`)
            } else {
                this.$message({
                    type: "success",
                    message: "插入成功！"
                })
                this.showAddFlowDlg = false
                await this.refresh()
            }
        }
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
