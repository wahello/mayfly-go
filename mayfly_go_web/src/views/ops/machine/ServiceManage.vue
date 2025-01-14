<template>
    <div class="file-manage">
        <el-dialog :title="title" v-model="visible" :destroy-on-close="true" :show-close="true" :before-close="handleClose" width="60%">
            <div class="toolbar">
                <div style="float: left">
                    <el-select v-model="type" @change="getScripts" size="mini" placeholder="请选择">
                        <el-option :key="0" label="私有" :value="0"> </el-option>
                        <el-option :key="1" label="公共" :value="1"> </el-option>
                    </el-select>
                </div>
                <div style="float: right">
                    <el-button @click="editScript(currentData)" :disabled="currentId == null" type="primary" icon="el-icon-tickets" size="mini" plain
                        >查看</el-button
                    >
                    <el-button v-auth="'machine:script:save'" type="primary" @click="editScript(null)" icon="el-icon-plus" size="mini" plain
                        >添加</el-button
                    >
                    <el-button
                        v-auth="'machine:script:del'"
                        :disabled="currentId == null"
                        type="danger"
                        @click="deleteRow(currentData)"
                        icon="el-icon-delete"
                        size="mini"
                        plain
                        >删除</el-button
                    >
                </div>
            </div>

            <el-table :data="scriptTable" @current-change="choose" stripe border size="mini" style="width: 100%">
                <el-table-column label="选择" width="55px">
                    <template #default="scope">
                        <el-radio v-model="currentId" :label="scope.row.id">
                            <i></i>
                        </el-radio>
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="名称" :min-width="70"> </el-table-column>
                <el-table-column prop="description" label="描述" :min-width="100" show-overflow-tooltip></el-table-column>
                <el-table-column prop="name" label="类型" :min-width="50">
                    <template #default="scope">
                        {{ enums.scriptTypeEnum.getLabelByValue(scope.row.type) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作">
                    <template #default="scope">
                        <el-button v-if="scope.row.id == null" @click="addFiles(scope.row)" type="success" icon="el-icon-success" size="mini" plain
                            >确定</el-button
                        >

                        <el-button
                            v-auth="'machine:script:run'"
                            v-if="scope.row.id != null"
                            @click="runScript(scope.row)"
                            type="primary"
                            icon="el-icon-tickets"
                            size="mini"
                            plain
                            >执行</el-button
                        >
                    </template>
                </el-table-column>
            </el-table>
        </el-dialog>

        <el-dialog title="脚本参数" v-model="scriptParamsDialog.visible" width="400px">
            <el-form ref="paramsForm" :model="scriptParamsDialog.params" label-width="70px" size="mini">
                <el-form-item v-for="item in scriptParamsDialog.paramsFormItem" :key="item.name" :prop="item.model" :label="item.name" required>
                    <el-input v-model="scriptParamsDialog.params[item.model]" :placeholder="item.placeholder" autocomplete="off"></el-input>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="hasParamsRun(currentData)" size="mini">确 定</el-button>
                </span>
            </template>
        </el-dialog>

        <el-dialog title="执行结果" v-model="resultDialog.visible" width="40%">
            <div style="white-space: pre-line; padding: 10px; color: #000000">
                {{ resultDialog.result }}
            </div>
        </el-dialog>

        <el-dialog
            v-if="terminalDialog.visible"
            title="终端"
            v-model="terminalDialog.visible"
            width="70%"
            :close-on-click-modal="false"
            :modal="false"
            @close="closeTermnial"
        >
            <ssh-terminal ref="terminal" :cmd="terminalDialog.cmd" :machineId="terminalDialog.machineId" height="600px" />
        </el-dialog>

        <script-edit
            v-model:visible="editDialog.visible"
            v-model:data="editDialog.data"
            :title="editDialog.title"
            v-model:machineId="editDialog.machineId"
            :isCommon="type == 1"
            @submitSuccess="submitSuccess"
        />
    </div>
</template>

<script lang="ts">
import { ref, toRefs, reactive, watch, defineComponent } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import SshTerminal from './SshTerminal.vue';
import { machineApi } from './api';
import enums from './enums';
import ScriptEdit from './ScriptEdit.vue';

export default defineComponent({
    name: 'ServiceManage',
    components: {
        ScriptEdit,
        SshTerminal,
    },
    props: {
        visible: { type: Boolean },
        machineId: { type: Number },
        title: { type: String },
    },
    setup(props: any, context) {
        const paramsForm: any = ref(null);
        const state = reactive({
            visible: false,
            type: 0,
            currentId: null,
            currentData: null,
            editDialog: {
                visible: false,
                data: null,
                title: '',
                machineId: 9999999,
            },
            scriptTable: [],
            scriptParamsDialog: {
                visible: false,
                params: {},
                paramsFormItem: [],
            },
            resultDialog: {
                visible: false,
                result: '',
            },
            terminalDialog: {
                visible: false,
                cmd: '',
                machineId: 0,
            },
        });

        watch(props, (newValue, oldValue) => {
            if (props.machineId) {
                getScripts();
            }
            state.visible = newValue.visible;
        });

        const getScripts = async () => {
            state.currentId = null;
            state.currentData = null;
            const machineId = state.type == 0 ? props.machineId : 9999999;
            const res = await machineApi.scripts.request({ machineId: machineId });
            state.scriptTable = res.list;
        };

        const runScript = async (script: any) => {
            // 如果存在参数，则弹窗输入参数后执行
            if (script.params) {
                state.scriptParamsDialog.paramsFormItem = JSON.parse(script.params);
                state.scriptParamsDialog.visible = true;
                return;
            }

            run(script);
        };

        // 有参数的脚本执行函数
        const hasParamsRun = async (script: any) => {
            // 如果脚本参数弹窗显示，则校验参数表单数据通过后执行
            if (state.scriptParamsDialog.visible) {
                paramsForm.value.validate((valid: any) => {
                    if (valid) {
                        run(script);
                        state.scriptParamsDialog.params = {};
                        state.scriptParamsDialog.visible = false;
                        paramsForm.value.resetFields();
                    } else {
                        return false;
                    }
                });
            }
        };

        const run = async (script: any) => {
            const noResult = script.type == enums.scriptTypeEnum['NO_RESULT'].value;
            // 如果脚本类型为有结果类型，则显示结果信息
            if (script.type == enums.scriptTypeEnum['RESULT'].value || noResult) {
                const res = await machineApi.runScript.request({
                    machineId: props.machineId,
                    scriptId: script.id,
                    params: state.scriptParamsDialog.params,
                });

                if (noResult) {
                    ElMessage.success('执行完成');
                    return;
                }
                state.resultDialog.result = res;
                state.resultDialog.visible = true;
                return;
            }

            if (script.type == enums.scriptTypeEnum['REAL_TIME'].value) {
                state.terminalDialog.cmd = script.script;
                state.terminalDialog.visible = true;
                state.terminalDialog.machineId = props.machineId;
                return;
            }
        };

        const closeTermnial = () => {
            state.terminalDialog.visible = false;
            state.terminalDialog.machineId = 0;
            // const t: any = this.$refs['terminal']
            // t.closeAll()
        };

        /**
         * 选择数据
         */
        const choose = (item: any) => {
            if (!item) {
                return;
            }
            state.currentId = item.id;
            state.currentData = item;
        };

        const editScript = (data: any) => {
            state.editDialog.machineId = props.machineId;
            state.editDialog.data = data;
            if (data) {
                state.editDialog.title = '查看编辑脚本';
            } else {
                state.editDialog.title = '新增脚本';
            }
            state.editDialog.visible = true;
        };

        const submitSuccess = () => {
            // this.delChoose()
            // this.search()
            getScripts();
        };

        const deleteRow = (row: any) => {
            ElMessageBox.confirm(`此操作将删除 [${row.name}], 是否继续?`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
            }).then(() => {
                machineApi.deleteScript
                    .request({
                        machineId: props.machineId,
                        scriptId: row.id,
                    })
                    .then((res) => {
                        getScripts();
                    });
                // 删除配置文件
            });
        };

        /**
         * 关闭取消按钮触发的事件
         */
        const handleClose = () => {
            context.emit('update:visible', false);
            context.emit('update:machineId', null);
            context.emit('cancel');
            state.scriptTable = [];
        };

        return {
            ...toRefs(state),
            paramsForm,
            enums,
            getScripts,
            runScript,
            hasParamsRun,
            closeTermnial,
            choose,
            editScript,
            submitSuccess,
            deleteRow,
            handleClose,
        };
    },
});
</script>
<style lang="sass">
</style>
