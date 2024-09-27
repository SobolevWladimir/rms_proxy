<template>
  <q-page class="q-pa-md">
    <q-list  bordered separator style="max-width: 318px">
      <q-item v-for="item in list" clickable v-ripple :key="item.id">
        <q-item-section>
          <q-item-label overline>{{item.url}}</q-item-label>
          <q-item-label>Логин: {{item.login}}</q-item-label>
          <q-item-label>Пароль: {{item.password}}</q-item-label>
        </q-item-section>
      </q-item>

      <q-item clickable v-ripple>
        <q-item-section>
          <q-item-label overline>url:</q-item-label>
          <q-item-label>login</q-item-label>
          <q-item-label>password</q-item-label>
        </q-item-section>
      </q-item>

      <q-item clickable v-ripple>
        <q-item-section>
          <q-item-label overline>url:</q-item-label>
          <q-item-label>login</q-item-label>
          <q-item-label>password</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
    <div class="q-pa-md q-gutter-sm">
    <q-btn color="primary" label="Добавить" />
  </div>
       <q-dialog v-model="showEdit" persistent>
      <q-card>
         <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Данные RMS</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-md">
            <q-input outlined v-model="editElem.url" label="Url" />
            <q-input outlined v-model="editElem.login" label="Логин" />
            <q-input outlined v-model="editElem.password" label="Пароль" />
          <q-checkbox v-model="editElem.needPassEncrupt" label="Преобразовать пароль в sha1" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Сохранить" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>

import { ref } from 'vue';

defineOptions({
  name: 'RmsPage',
});
const showEdit = ref(true);
const list = ref([
  {
    id: '234342',
    url: 'http://192.168.0.155:9080',
    login: 'adm',
    password: '123',
    needPassEncrupt: true,
  },
]);
const editElem = ref({
  id: null,
  url: '',
  login: '',
  password: '',
  needPassEncrupt: false,
});

</script>
