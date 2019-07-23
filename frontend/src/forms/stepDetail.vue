<template>
<el-form ref="form" label-width="80px">
    <el-form-item label="操作标识">
        {{selStep.operKey}}
    </el-form-item>
    <el-form-item label="依赖" v-show="selStep.requires && selStep.requires.length !== 0">
        <el-tag v-for="require in selStep.requires" :key="require" :closable="mode !== 'display'">
            {{require}}
        </el-tag>
    </el-form-item>
    <el-form-item label="备注">
        <el-input v-if="mode !== 'display'" v-model="selStep.desc"/>
        <p v-else>{{selStep.desc}}</p>
    </el-form-item>
    <el-form-item label="输入" v-show="mode === 'editing-flow' && selStep.inputs && selStep.inputs.length !== 0">
        <div class="card input-card" v-for="(input, symbol) in selStep.inputs" :key="symbol">
            <div class="card-body input-card-body">
                <el-row>
                    <el-col :span="11">{{symbol}}</el-col>
                    <el-col :span="2">
                        <i class="el-icon-arrow-right"></i>
                    </el-col>
                    <el-col :span="11">
                        <el-dropdown class="float-right" trigger="click" @command="hdlSelInput">
                            <span class="el-dropdown-link">
                                {{input === "" ? "选择局部变量" : input}}<i class="el-icon-arrow-down el-icon--right"></i>
                            </span>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item :command="`${symbol}:ctx`">context</el-dropdown-item>
                                <el-dropdown-item v-for="lv in locVars" :key="lv" :command="`${symbol}:${lv}`">{{lv}}</el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                    </el-col>
                </el-row>
            </div>
        </div>
    </el-form-item>
    <el-form-item label="输出" v-show="mode === 'editing-flow' && selStep.outputs && selStep.outputs.length !== 0">
        <el-tag v-for="output in selStep.outputs" :key="output" type="success">{{output}}</el-tag>
    </el-form-item>
    <el-form :inline="true" label-width="80px" v-show="mode === 'editing-flow'">
        <el-form-item label="标识">
            <el-select v-model="selStep.symbol" placeholder="请选择">
                <el-option v-for="(value, name) in spcSymbols" :key="name" :label="name" :value="value"/>
            </el-select>
        </el-form-item>
        <el-form-item label="效果">
            <el-tag>abcd</el-tag>
        </el-form-item>
    </el-form>
    <el-form-item label="进退格" v-show="mode === 'new'">
        <el-checkbox v-model="enableBlk">启用进退格</el-checkbox>
        <el-switch v-model="blkInOut" active-text="进格" inactive-text="退格" :disabled="!enableBlk" @change="hdlSwhBlkInOut"/>
    </el-form-item>
    <el-form-item label="代码">
        <el-input type="textarea" v-model="selStep.code" :autosize="{minRows: 5, maxRows: 100}" :disabled="mode === 'display'"/>
    </el-form-item>
</el-form>
</template>

<script>
import backend from "../async/backend"

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
            blkInOut: true
        }
    },
    async created() {
        if (this.preMode && this.preMode.length !== 0) {
            this.mode = this.preMode
        }
        let res = await backend.specials()
        if (typeof res === "string") {
            this.$message.error(`查询特殊标识失败：${res}`)
        } else {
            this.spcSymbols = res.data.data.values || {}
        }
    },
    methods: {
        hdlSelInput(cmd) {
            let kvs = cmd.split(":")
            this.selStep.inputs[kvs[0]] = kvs[1]
        },
        hdlSwhBlkInOut() {
            if (this.enableBlk) {
                if (this.mode === "editing-flow") {
                    // 使用模板，blockIn在new模式下已经确定，无法在editing-flow修改
                    if (!this.selStep.blockIn) {
                        this.blkInOut = false
                        this.selStep.blockOut = true
                    }
                } else {
                    this.selStep.blockIn = this.blkInOut
                    this.selStep.blockOut = !this.blkInOut
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
