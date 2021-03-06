<template>
<el-form ref="form" label-width="80px">
    <el-form-item label="操作标识">{{selStep.key}}</el-form-item>
    <el-form ref="new-require-form" v-model="newRequire" v-show="mode === 'add-temp-step'" label-width="80px">
        <el-form-item label="添加依赖" :rules="[
            { required: true, message: '请输入依赖路径', trigger: 'blur' }
        ]">
            <el-row class="m-0" :gutter="5">
                <el-col class="p-0" :span="18">
                    <el-input v-model="newRequire.input" placeholder="请输入依赖路径"/>
                </el-col>
                <el-col class="p-0" :span="6">
                    <el-button class="float-right" @click="addRequire">添加</el-button>
                </el-col>
            </el-row>
        </el-form-item>
    </el-form>
    <el-form-item label="依赖" v-show="selStep.requires && selStep.requires.length !== 0">
        <el-tag v-for="require in selStep.requires" :key="require" :closable="mode !== 'display'">
            {{require}}
        </el-tag>
    </el-form-item>
    <el-form-item label="备注">
        <el-input v-if="mode !== 'display'" v-model="selStep.desc"/>
        <p v-else>{{selStep.desc}}</p>
    </el-form-item>
    <el-form-item label="输入" v-show="mode === 'editing-step' && selStep.inputs && selStep.inputs.length !== 0">
        <div class="card input-card" v-for="(input, symbol) in selStep.inputs" :key="symbol">
            <div class="card-body input-card-body">
                <el-row class="plr-10">
                    <el-col :span="11">{{symbol}}：{{input.desc}}</el-col>
                    <el-col :span="2">
                        <i class="el-icon-arrow-left"></i>
                    </el-col>
                    <el-col :span="11">
                        <el-dropdown v-if="symbol !== inputTextVar" class="float-right" trigger="click" @command="hdlSelInput">
                            <span class="el-dropdown-link">
                                {{!input.name ? "选择局部变量" : input.name}}<i class="el-icon-arrow-down el-icon--right"></i>
                            </span>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item v-for="lv in locVars" :key="lv" :command="`${symbol}:${lv}`">{{lv}}</el-dropdown-item>
                                <el-dropdown-item :command="symbol">文本输入</el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                        <el-input v-else v-model="selStep.inputs[symbol].name">
                            <el-button slot="prepend" icon="el-icon-check" @click="hdlInputVar(true)"/>
                            <el-button slot="append" icon="el-icon-close" @click="hdlInputVar(false)"/>
                        </el-input>
                    </el-col>
                </el-row>
            </div>
        </div>
    </el-form-item>
    <el-form-item label="输出" v-show="mode === 'editing-step' && selStep.outputs && selStep.outputs.length !== 0">
        <el-tag v-for="output in selStep.outputs" :key="output" type="success">{{output}}</el-tag>
    </el-form-item>
    <el-form :inline="true" label-width="80px" v-show="mode === 'editing-step' && false">
        <el-form-item label="标识">
            <el-select v-model="selStep.symbol" placeholder="请选择">
                <el-option v-for="(value, name) in spcSymbols" :key="name" :label="name" :value="value"/>
            </el-select>
        </el-form-item>
        <el-form-item label="效果">
            <el-tag>abcd</el-tag>
        </el-form-item>
    </el-form>
    <el-form-item label="代码">
        <el-input type="textarea" v-model="selStep.code" :autosize="{minRows: 5, maxRows: 100}" :disabled="mode === 'display'"/>
    </el-form-item>
</el-form>
</template>

<script>
import backend from "../backend"

export default {
    props: {
        "selStep": Object,
        "preMode": String,
        "locVars": Array
    },
    data() {
        return {
            mode: "display",
            spcSymbols: {},
            enableBlk: false,
            blkInOut: true,
            newRequire: {input: ""},
            inputTextVar: "",
            procCode: ""
        }
    },
    async created() {
        if (this.preMode && this.preMode.length !== 0) {
            this.mode = this.preMode
        }
        let res = await backend.qryStepSymbols()
        if (typeof res === "string") {
            this.$message.error(`查询特殊标识失败：${res}`)
        } else {
            this.spcSymbols = res.values || {}
        }
    },
    methods: {
        hdlSelInput(cmd) {
            let kvs = cmd.split(":")
            if (kvs.length !== 2) {
                this.inputTextVar = cmd
            } else {
                this.selStep.inputs[kvs[0]].name = kvs[1]
            }
        },
        hdlSwhBlkInOut() {
            if (this.enableBlk) {
                if (this.mode === "editing-step") {
                    // 使用模板，blockIn在new模式下已经确定，无法在editing-step修改
                    if (!this.selStep.blockIn) {
                        this.blkInOut = false
                        this.selStep.blockOut = true
                    }
                } else {
                    this.selStep.blockIn = this.blkInOut
                    this.selStep.blockOut = !this.blkInOut
                }
            }
        },
        addRequire() {
            this.$refs["new-require-form"].validate(valid => {
                if (!valid) {
                    return false
                }
            })
        },
        hdlInputVar(useInput) {
            if (!useInput) {
                this.selStep.inputs[this.inputTextVar].name = ""
            } else {
                // 更新备注
                // this.selStep.desc = this.selStep.desc.replace(
                //     `%${this.inputTextVar}%`,
                //     this.selStep.inputs[this.inputTextVar].name
                // )
                // 更新输出
                for (let idx in this.selStep.outputs) {
                    this.selStep.outputs[idx] = _.replace(
                        this.selStep.outputs[idx],
                        `%${this.inputTextVar}%`,
                        this.selStep.inputs[this.inputTextVar].name
                    )
                }
            }
            this.inputTextVar = ""
        },
        toProcCode() {
            this.procCode = this.selStep.code
            for (let symbol in this.selStep.inputs) {
                const input = this.selStep.inputs[symbol]
                if (input.length !== 0) {
                    this.procCode = _.replace(this.procCode, `%${symbol}%`, input)
                }
            }
        }
    }
}
</script>

<style lang="scss">
.input-card {
    margin-bottom: 1vh;
}
.input-card:last-child {
    margin-bottom: 0 !important;
}
.input-card-body {
    padding: 0.2vh 0.5vw;
    .el-dropdown {
        cursor: pointer;
    }
}
.highlight {
    background-color: red;
}
</style>
