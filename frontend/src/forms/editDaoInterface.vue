<template>
<el-form :model="daoInterface" label-width="80px" ref="edit-dao-interface-form">
    <el-form-item label="接口名" :rules="[{
        required: true, message: '需要指定接口名', trigger: 'blur'
    }]" prop="name">
        <el-input placeholder="输入接口名字" v-model="daoInterface.name"/>
    </el-form-item>
    <el-form :model="newParam" label-width="80px" ref="add-param-form">
        <el-form-item label="添加参数" :rules="[
            { required: true, message: '需要指定参数名', trigger: 'blur' },
            { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
        ]" prop="name">
            <el-input placeholder="输入参数名" v-model="newParam.name"/>
        </el-form-item>
        <el-form-item :rules="[
            { required: true, message: '请选择参数类型', trigger: 'change' }
        ]" prop="type">
            <el-row class="m-0" :gutter="5">
                <el-col class="p-0" :span="18">
                    <el-select v-model="newParam.type" style="width: 100%" placeholder="请输入参数类型">
                        <el-option v-for="option in supportTypes" :key="option.title" :label="option.title" :value="option.value"/>
                    </el-select>
                </el-col>
                <el-col class="p-0" :span="6">
                    <el-button class="float-right" @click="addParam">添加</el-button>
                </el-col>
            </el-row>
        </el-form-item>
    </el-form>
    <el-form-item label="参数" prop="params" v-show="daoInterface.params.length !== 0">
        <el-tag v-for="param in daoInterface.params" :key="param.name" closable @close="delParam(param.name)">{{param.name}}</el-tag>
    </el-form-item>
    <el-form :model="newReturn" label-width="80px" ref="add-return-form">
        <el-form-item label="添加返回" :rule="[
            { required: true, message: '请选择返回类型', trigger: 'change' }
        ]" prop="type">
            <el-row class="m-0" :gutter="5">
                <el-col class="p-0" :span="18">
                    <el-select v-model="newReturn.type" style="width: 100%" placeholder="请输入返回值类型">
                        <el-option v-for="option in supportTypes" :key="option.title" :label="option.title" :value="option.value"/>
                    </el-select>
                </el-col>
                <el-col class="p-0" :span="6">
                    <el-button class="float-right" @click="addReturn">添加</el-button>
                </el-col>
            </el-row>
        </el-form-item>
    </el-form>
    <el-form-item label="返回值" prop="returns" v-show="daoInterface.returns.length !== 0">
        <el-tag v-for="(ret, idx) in daoInterface.returns" :key="idx" closable @close="delReturn(idx)">{{ret}}</el-tag>
    </el-form-item>
    <el-form :model="newRequire" label-width="80px" ref="add-require-form">
        <el-form-item label="添加依赖" :rule="[
            { required: true, message: '请选择依赖模块', trigger: 'change' }
        ]" prop="mod">
            <el-row class="m-0" :gutter="5">
                <el-col class="p-0" :span="18">
                    <el-select v-model="newRequire.mod" style="width: 100%" placeholder="请输入依赖模块">
                        <el-option v-for="option in supportTypes" :key="option.title" :label="option.title" :value="option.value"/>
                    </el-select>
                </el-col>
                <el-col class="p-0" :span="6">
                    <el-button class="float-right" @click="addRequire">添加</el-button>
                </el-col>
            </el-row>
        </el-form-item>
    </el-form>
    <el-form-item label="依赖模块" prop="requires" v-show="daoInterface.requires.length !== 0">
        
    </el-form-item>
    <el-form-item label="描述" prop="desc">
        <el-input type="textarea" :rows="2" placeholder="请输入描述" v-model="daoInterface.desc"/>
    </el-form-item>
</el-form>
</template>

<script>
import utils from "../utils"

export default {
    data() {
        return {
            supportTypes: utils.supportTypes,

            daoInterface: {
                name: "",
                params: [],
                returns: [],
                requires: []
            },
            newParam: {
                name: "",
                type: ""
            },
            newReturn: { type: "" },
            newRequire: { mod: "" }
        }
    },
    async created() {
        await this.refsRequires()
    },
    methods: {
        addParam() {
            this.$refs["add-param-form"].validate(valid => {
                if (!valid) {
                    return false
                }
                this.daoInterface.params.push({
                    name: this.newParam.name,
                    type: this.newParam.type
                })
                this.newParam.name = ""
                this.newParam.type = ""
            })
        },
        delParam(pname) {
            this.daoInterface.params.pop(ele => ele.name === pname)
        },
        addReturn() {
            this.$refs["add-return-form"].validate(valid => {
                if (!valid) {
                    return false
                }
                this.daoInterface.returns.push(this.newReturn.type)
                this.newReturn.type = ""
            })
        },
        delReturn(rindex) {
            this.daoInterface.returns.splice(rindex, 1)
        },
        async refsRequires() {

        },
        addRequire() {

        }
    }
}
</script>