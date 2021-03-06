<template>
<dashboard>
    <info-bar @select-api="selectApi" @add-api="addApi" @del-api="delApi"/>
    <div v-if="forceUpdateFlg">
        <div id="pnlFlows" class="w-100 h-100" style="position: absolute; padding-bottom: 50px">
            <step-block
                v-for="(step, index) in selApi.steps" :key="index"
                :step="step"
                @show-detail="showOperDetail"
                @be-deleted="refresh"
            />
            <el-button
                type="success" icon="el-icon-plus" circle
                v-for="btn in istStepBtns" :key="btn.nsuffix"
                :name="`istStepBtn${btn.nsuffix}`"
                class="ml-0" style="position:absolute"
                @click="insertStep(btn)"
            />
        </div>
        <svg id="pnlGraphs" class="w-100" style="position: absolute; z-index: -100; height: 100%">
            <step-link v-for="btn in istStepBtns" :key="btn.nsuffix" :istStepBtn="btn"/>
        </svg>
    </div>
    <el-button id="btnApiInfo" type="primary" size="small" icon="el-icon-info" v-if="selApi" @click="showApiInfo">
        {{selApi.name}}
    </el-button>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog :title="`步骤 #${selStep.index + 1}`" :visible.sync="showStepDtlDlg" :modal-append-to-body="false" width="50vw">
        <step-detail ref="step-detail-form" :selStep="selStep"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="chgOperStep">编 辑</el-button>
            <el-button type="primary">确 定</el-button>
        </div>
    </el-dialog>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog :title="`添加步骤 #${clkIstStepBtn.prev ? clkIstStepBtn.prev.index + 2 : 1}`" :visible.sync="showAddStepDlg" :modal-append-to-body="false" width="50vw">
        <edit-step ref="add-step-form" :locVars="clkIstStepBtn.locVars"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddStepDlg = false">取 消</el-button>
            <el-button type="primary" @click="addStep">确 定</el-button>
        </div>
    </el-dialog>
</dashboard>
</template>

<script>
import _ from "lodash"

import dashboard from "../layouts/dashboard"
import infoBar from "../components/infoBar"
import stepDetail from "../forms/stepDetail"
import backend from "../backend"
import apiInfo from "../forms/apiInfo"
import editStep from "../forms/editStep"
import stepBlock from "../components/stepBlock"
import stepLink from "../components/stepLink"

export default {
    components: {
        "dashboard": dashboard,
        "info-bar": infoBar,
        "step-detail": stepDetail,
        "edit-step": editStep,
        "step-block": stepBlock,
        "step-link": stepLink
    },
    data() {
        return {
            forceUpdateFlg: true,
            selApi: {steps: []},
            selStep: {index: -1},
            index: 1,
            showStepDtlDlg: false,
            showAddStepDlg: false,
            istStepBtns: [],
            clkIstStepBtn: {prev: null}
        }
    },
    updated() {
        let pnlFlows = document.getElementById("pnlFlows")
        if (pnlFlows) {
            d3.select("#pnlGraphs")
            .style("height", `${pnlFlows.scrollHeight}px`)
        }
    },
    methods: {
        async refresh() {
            let res = await backend.qryApiByName(this.selApi.name)
            if (typeof res === "string") {
                this.$message.error(`查询名为${this.selApi.name}的接口时发生错误：${res}`)
            } else {
                this.selectApi(res)
            }
        },
        async selectApi(selApi) {
            this.selApi = selApi
            this.istStepBtns = []
            let locVars = this.selApi.params ? _.keys(this.selApi.params) : []
            this.selApi.steps = this.selApi.steps ? this.selApi.steps.map((step, idx) => {
                // 做一些处理：只包含一个元素的输出数组全部设为空
                if (!step.outputs || (step.outputs.length === 1 && step.outputs[0] === "")) {
                    step.outputs = []
                }
                if (!step.requires || (step.requires.length === 1 && step.requires[0] === "")) {
                    step.requires = []
                }
                step.index = idx

                if (idx === this.selApi.steps.length - 1) {
                    // 如果最后一个步骤的标识不是结尾标识，则添加按钮用于后续增加步骤
                    if (!step.symbol || step.symbol & 4 /* SpcSymbol_END */ === 0) {
                        this.istStepBtns.push({
                            apiName: selApi.name,
                            nsuffix: `_${idx}`,
                            prev: step,
                            next: null,
                            locVars
                        })
                    } else {
                        step.isLast = true
                    }
                } else {
                    this.istStepBtns.push({
                        apiName: selApi.name,
                        nsuffix: `_${idx}_${idx + 1}`,
                        prev: step,
                        next: this.selApi.steps[idx + 1],
                        locVars
                    })
                }

                if (idx !== 0) {
                    locVars.concat(step.outputs)
                }
                return step
            }) : []
            if (this.selApi.name && this.istStepBtns.length === 0) {
                // 没有按钮，说明流程中没有步骤，添加一个按钮用于初始化
                this.istStepBtns.push({
                    apiName: selApi.name,
                    nsuffix: "__0",
                    locVars
                })
            }
            // 强制pnlFlows刷新
            // NOTE: 如果不强制刷新，同名的step块会相互覆盖并影响link的定位
            this.forceUpdateFlg = false
            this.$nextTick(() => {
                this.forceUpdateFlg = true
            })
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
                showConfirmButton: false,
                customClass: "w-50"
            }).catch(err => {})
        },
        addApi(newApi) {
            console.log(newApi)
        },
        delApi() {
            this.selectApi({steps: []})
        },
        async addStep() {
            let form = this.$refs["add-step-form"]
            let res = await backend.addStep({
                index: this.clkIstStepBtn.prev ? this.clkIstStepBtn.prev.index + 1 : 0,
                step: Object.assign(form.step, {
                    apiName: this.selApi.name
                })
            })
            if (typeof res === "string") {
                this.$message.error(`插入流程发生错误：${res}`)
            } else {
                this.$message({
                    type: "success",
                    message: "插入成功！"
                })
                this.showAddStepDlg = false
                await this.refresh()
            }
        },
        showOperDetail(step) {
            this.selStep = step
            this.showStepDtlDlg = true
        },
        insertStep(istBtn) {
            this.clkIstStepBtn = istBtn
            this.showAddStepDlg = true
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
    right: 20px;
    bottom: 0;
    border-radius: 4px 4px 0 0;
}
</style>
