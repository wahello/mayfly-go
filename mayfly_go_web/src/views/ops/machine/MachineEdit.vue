<template>
    <div>
        <el-dialog :title="title" v-model="visible" :show-close="false" :before-close="cancel" width="35%">
            <el-form :model="form" ref="machineForm" :rules="rules" label-width="85px" size="small">
                <!-- <el-form-item prop="projectId" label="项目:" required>
                    <el-select style="width: 100%" v-model="form.projectId" placeholder="请选择项目" @change="changeProject" filterable>
                        <el-option v-for="item in projects" :key="item.id" :label="`${item.name} [${item.remark}]`" :value="item.id"> </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item prop="envId" label="环境:" required>
                    <el-select @change="changeEnv" style="width: 100%" v-model="form.envId" placeholder="请选择环境">
                        <el-option v-for="item in envs" :key="item.id" :label="`${item.name} [${item.remark}]`" :value="item.id"> </el-option>
                    </el-select>
                </el-form-item> -->
                <el-form-item prop="name" label="名称:" required>
                    <el-input v-model.trim="form.name" placeholder="请输入机器别名" auto-complete="off"></el-input>
                </el-form-item>
                <el-form-item prop="ip" label="ip:" required>
                    <el-input v-model.trim="form.ip" placeholder="请输入主机ip" auto-complete="off"></el-input>
                </el-form-item>
                <el-form-item prop="port" label="port:" required>
                    <el-input type="number" v-model.trim="form.port" placeholder="请输入端口"></el-input>
                </el-form-item>
                <el-form-item prop="username" label="用户名:" required>
                    <el-input v-model.trim="form.username" placeholder="请输入用户名"></el-input>
                </el-form-item>
                <el-form-item prop="password" label="密码:" required>
                    <el-input
                        type="password"
                        show-password
                        v-model.trim="form.password"
                        placeholder="请输入密码"
                        autocomplete="new-password"
                    ></el-input>
                </el-form-item>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button type="primary" :loading="btnLoading" @click="btnOk" size="mini">确 定</el-button>
                    <el-button @click="cancel()" size="mini">取 消</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script lang="ts">
import { toRefs, reactive, watch, onMounted, defineComponent, ref } from 'vue';
import { machineApi } from './api';
import { projectApi } from '../project/api.ts';
import { ElMessage } from 'element-plus';

export default defineComponent({
    name: 'MachineEdit',
    props: {
        visible: {
            type: Boolean,
        },
        machine: {
            type: [Boolean, Object],
        },
        title: {
            type: String,
        },
    },
    setup(props: any, { emit }) {
        const machineForm: any = ref(null);
        const state = reactive({
            visible: false,
            form: {
                id: null,
                name: null,
                port: 22,
                username: null,
                password: null,
            },
            btnLoading: false,
            rules: {
                projectId: [
                    {
                        required: true,
                        message: '请选择项目',
                        trigger: ['change', 'blur'],
                    },
                ],
                envId: [
                    {
                        required: true,
                        message: '请选择环境',
                        trigger: ['change', 'blur'],
                    },
                ],
                name: [
                    {
                        required: true,
                        message: '请输入别名',
                        trigger: ['change', 'blur'],
                    },
                ],
                ip: [
                    {
                        required: true,
                        message: '请输入主机ip',
                        trigger: ['change', 'blur'],
                    },
                ],
                port: [
                    {
                        required: true,
                        message: '请输入端口',
                        trigger: ['change', 'blur'],
                    },
                ],
                username: [
                    {
                        required: true,
                        message: '请输入用户名',
                        trigger: ['change', 'blur'],
                    },
                ],
                password: [
                    {
                        required: true,
                        message: '请输入密码',
                        trigger: ['change', 'blur'],
                    },
                ],
            },
        });

        watch(props, async (newValue, oldValue) => {
            state.visible = newValue.visible;
            if (newValue.machine) {
                state.form = { ...newValue.machine };
            } else {
                state.form = { port: 22 } as any;
            }
        });

        const btnOk = async () => {
            machineForm.value.validate((valid: boolean) => {
                if (valid) {
                    machineApi.saveMachine.request(state.form).then((res: any) => {
                        ElMessage.success('保存成功');
                        emit('val-change', state.form);
                        state.btnLoading = true;
                        setTimeout(() => {
                            state.btnLoading = false;
                        }, 1000);

                        cancel();
                    });
                } else {
                    ElMessage.error('请正确填写信息');
                    return false;
                }
            });
        };

        const cancel = () => {
            emit('update:visible', false);
            emit('cancel');
            setTimeout(() => {
                machineForm.value.resetFields();
                //  重置对象属性为null
                state.form = {} as any;
            }, 200);
        };

        return {
            ...toRefs(state),
            machineForm,
            btnOk,
            cancel,
        };
    },
});
</script>
<style lang="scss">
</style>
