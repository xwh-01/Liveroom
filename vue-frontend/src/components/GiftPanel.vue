<template>
  <div class="gift-buttons">
    <el-button
      type="danger"
      size="small"
      :disabled="!connected"
      @click="$emit('send-gift', 'heart')"
      plain
    >
      ❤ 小心心 (10分)
    </el-button>
    <el-button
      type="warning"
      size="small"
      :disabled="!connected"
      @click="$emit('send-gift', 'rocket')"
      plain
    >
      🚀 火箭 (100分)
    </el-button>
  </div>
</template>

<script setup>
defineProps({
  connected: { type: Boolean, default: false }
})
defineEmits(['send-gift'])
</script>

<style scoped>
.gift-buttons {
  display: flex;
  gap: 8px;
}
</style>
