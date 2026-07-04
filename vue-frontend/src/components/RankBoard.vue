<template>
  <div class="card-panel rank-board">
    <div class="panel-title">🎖 礼物排行榜 TOP10</div>
    <div v-if="rankings.length === 0" class="rank-empty">暂无排行</div>
    <div class="rank-list">
      <div
        v-for="(item, idx) in rankings"
        :key="idx"
        class="rank-item"
        :class="{ 'top-3': idx < 3 }"
      >
        <span class="rank-index" :class="'rank-' + (idx + 1)">{{ idx + 1 }}</span>
        <span class="rank-user">{{ item.user_id }}</span>
        <span class="rank-score">{{ item.score }}分</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  rankings: { type: Array, default: () => [] }
})
</script>

<style scoped>
.rank-empty {
  color: #c0c4cc;
  text-align: center;
  margin: 20px 0;
  font-size: 13px;
}

.rank-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f5f5f5;
  font-size: 13px;
}

.rank-item.top-3 {
  font-weight: 600;
}

.rank-index {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  background: #ebeef5;
  color: #909399;
  margin-right: 10px;
  flex-shrink: 0;
}

.rank-1 { background: #ff4757; color: #fff; }
.rank-2 { background: #ff6b81; color: #fff; }
.rank-3 { background: #ffa502; color: #fff; }

.rank-user {
  flex: 1;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-score {
  color: #e6a23c;
  font-weight: 600;
  flex-shrink: 0;
}
</style>
