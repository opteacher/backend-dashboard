<template>
<el-form :model="model" label-width="80px" label-position="right" ref="edit-model-form">
    <el-form-item :label="`${structFlag ? '结构名' : '模块名'}`" :rules="[
        { required: true, message: '请输入模块名称', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
    ]" prop="name">
        <el-input v-model="model.name"/>
    </el-form-item>
    <el-form-item label="属性">
        <el-row class="m-0" :gutter="5">
            <el-col class="p-0" :span="18">
                <el-form-item :rules="[
                    { required: true, message: '请输入属性名称', trigger: 'blur' },
                    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
                ]" prop="propName">
                    <el-input v-model="model.propName"/>
                </el-form-item>
            </el-col>
            <el-col class="p-0" :span="6">
                <el-button class="float-right" @click="addProp">添加</el-button>
            </el-col>
        </el-row>
    </el-form-item>
    <el-form-item v-for="(prop, index) in model.props" :key="prop.name" :label="prop.name">
        <el-row class="m-0" :gutter="5">
            <el-col class="p-0" :span="18">
                <el-form-item :rules="[
                    { required: true, message: '请选择属性类型', trigger: 'change' }
                ]" :prop="`props.${index}.type`">
                    <el-select v-model="prop.type" style="width: 100%">
                        <el-option v-for="option in supportTypes" :key="option.title" :label="option.title" :value="option.value"/>
                    </el-select>
                </el-form-item>
            </el-col>
            <el-col class="p-0" :span="6">
                <el-button class="float-right" @click.prevent="deleteProp(prop.name)">删除</el-button>
            </el-col>
        </el-row>
    </el-form-item>
    <el-form-item label="RPC接口" v-if="!structFlag && persistDaoGrps.length !== 0">
        <el-col :span="18">
            <el-checkbox-group class="w-100" v-if="selPersistDao" v-model="model.methods">
                <el-checkbox-button v-if="selPersistDao['insert']" label="insert" name="methods">增</el-checkbox-button>
                <el-checkbox-button v-if="selPersistDao['delete']" label="delete" name="methods">删</el-checkbox-button>
                <el-checkbox-button v-if="selPersistDao['update']" label="update" name="methods">改</el-checkbox-button>
                <el-checkbox-button v-if="selPersistDao['query']" label="query" name="methods">查</el-checkbox-button>
                <el-checkbox-button v-if="selPersistDao['quertAll']" label="queryAll" name="methods">全查</el-checkbox-button>
            </el-checkbox-group>
        </el-col>
        <el-col :span="6">
            <el-select class="w-100" v-model="selPersistDao">
                <el-option v-for="(group, gname) in persistDaoGrps" :key="gname" :label="gname" :value="group"/>
            </el-select>
        </el-col>
    </el-form-item>
</el-form>
</template>

<script>
import utils from "../utils"
import backend from '../backend'

export default {
    props: {
        "input-model": Object,
        "structFlag": Boolean
    },
    data() {
        return {
            supportTypes: utils.supportTypes,
            model: {
                id: "",
                name: "",
                props: [],
                propName: "",
                methods: [],
                x: 0,
                y: 0,
                width: 400,
                height: 300
            },
            selPersistDao: null,
            persistDaoGrps: {}
        }
    },
    async created() {
        if (this["input-model"] && this["input-model"].name) {
            this.model.name = this["input-model"].name
        }
        if (this["input-model"] && this["input-model"].props) {
            this.model.props = this["input-model"].props
        }
        let res = await backend.qryTempApisByCategory("persist")
        if (typeof res === "string") {
            this.$message.error(`查询持久化DAO时发生错误：${res}`)
        } else if (!res.infos) {
            return
        } else {
            let categories = {"": []}
            let firstGroup = ""
            for (let info of res.infos) {
                info.symbol = info.symbol || ""
                if (categories[info.group]) {
                    categories[info.group][info.symbol] = info
                } else {
                    categories[info.group] = {[info.symbol]: info}
                    if (firstGroup.length === 0) {
                        firstGroup = info.group
                    }
                }
            }
            if (categories[""].length === 0) {
                delete categories[""]
            }
            this.persistDaoGrps = categories
            if (firstGroup.length !== 0) {
                this.selPersistDao = categories[firstGroup]
            }
        }
    },
    methods: {
        addProp() {
            this.$refs["edit-model-form"].validate(valid => {
                if (valid) {
                    this.model.props.push({
                        name: this.model.propName
                    })
                    this.model.propName = ""
                } else {
                    return false
                }
            })
        },
        deleteProp(propName) {
            this.model.props.pop(ele => ele.name === propName)
        },
        resetModel() {
            this.model = {
                id: "",
                name: "",
                props: [],
                propName: "",
                methods: [],
                x: 0,
                y: 0,
                width: 400,
                height: 300
            }
        }
    }
}
</script>
