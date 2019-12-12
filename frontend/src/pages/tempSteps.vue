<template>
<dashboard>
    <div class="table-container">
        <el-button type="primary" @click="showAddTempStep = true">添加模板步骤</el-button>
        <el-table class="mt-10" :data="tempSteps" style="width: 100%">
            <el-table-column prop="idenKey" label="标识"/>
            <el-table-column prop="group" label="组"/>
            <el-table-column prop="requires" label="依赖"/>
            <el-table-column prop="desc" label="描述"/>
            <el-table-column prop="inputs" label="输入（槽）">
                <template slot-scope="scope">
                    <ul class="list-unstyled">
                        <li v-for="(input, slot) in scope.row.inputs" :key="slot">{{slot}}</li>
                    </ul>
                </template>
            </el-table-column>
            <el-table-column prop="outputs" label="输出（变量）">
                <template slot-scope="scope">
                    <ul class="list-unstyled">
                        <li v-for="output in scope.row.outputs" :key="output">{{output}}</li>
                    </ul>
                </template>
            </el-table-column>
            <el-table-column prop="code" label="代码">
                <template slot-scope="scope">
                    <el-button size="small" @click="showCode(scope.row.code)">查看代码</el-button>
                </template>
            </el-table-column>
        </el-table>
    </div>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog title="添加模板步骤" :visible.sync="showAddTempStep" :modal-append-to-body="false" width="50vw">
        <step-detail ref="add-temp-step-form" :selStep="addStep" preMode="add-temp-step" :stepInfo="{locVars:[]}"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddTempStep = false">取 消</el-button>
            <el-button type="primary" @click="addTempStep">确 定</el-button>
        </div>
    </el-dialog>
</dashboard>
</template>

<script>
import dashboard from "../layouts/dashboard"
import backend from '../backend'
import codeView from "../forms/codeView"
import stepDetail from "../forms/stepDetail"

export default {
    components: {
        "dashboard": dashboard,
        "step-detail": stepDetail
    },
    data() {
        return {
            showAddTempStep: false,
            tempSteps: [],
            addStep: {}
        }
    },
    async created() {
        await this.refresh()
    },
    methods: {
        async refresh() {
            let res = await backend.qryAllTempStep()
            if (typeof res === "string") {
                this.$message.error(`查询模板步骤发生错误：${res}`)
            } else {
                this.tempSteps = res.steps
            }
        },
        addTempStep() {
            
        },
        showCode(code) {
            this.$msgbox({
                title: "查看代码",
                message: this.$createElement(codeView, {
                    props: {code}
                }),
                showConfirmButton: false
            }).catch(err => {})
        }
    }
}
</script>