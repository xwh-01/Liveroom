<template>
  <div class="card-panel live-player">
    <div class="panel-title">直播画面</div>
    <div class="player-area">
      <div class="player-badge">LIVE</div>
      <div class="player-room">Live Room {{ defaultRoom }}</div>
      <div class="player-status">主播正在直播中...</div>
      <div class="player-wave">
        <span class="bar"></span>
        <span class="bar"></span>
        <span class="bar"></span>
        <span class="bar"></span>
        <span class="bar"></span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  defaultRoom: { type: String, default: '1001' }
})
</script>

<style scoped>
.player-area {
  height: 260px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.player-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  background: #ff4757;
  color: #fff;
  padding: 2px 10px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 1px;
  animation: pulse 1.5s infinite;
}

.player-room {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.player-status {
  color: #a4b0be;
  font-size: 13px;
  margin-bottom: 24px;
}

.player-wave {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 30px;
}

.player-wave .bar {
  width: 4px;
  background: #667eea;
  border-radius: 2px;
  animation: wave 1.2s ease-in-out infinite;
}

.player-wave .bar:nth-child(2) { animation-delay: 0.1s; height: 18px; }
.player-wave .bar:nth-child(3) { animation-delay: 0.2s; height: 28px; }
.player-wave .bar:nth-child(4) { animation-delay: 0.3s; height: 14px; }
.player-wave .bar:nth-child(5) { animation-delay: 0.4s; height: 22px; }

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

@keyframes wave {
  0%, 100% { height: 8px; }
  50% { height: 24px; }
}
</style>
