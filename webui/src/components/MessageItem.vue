<script>
export default {
    props: ['messageId', 'author', 'text', 'hours']
}
</script>

<template>
  <div class="messageItem">
    <div class="messageBox">
      <div class="messageHeader">
        <strong>{{ author }}</strong>
        <button class="optionsButton">⋮</button>
      </div>
      <div class="messageText">{{ text }}</div>
      <div class="messageTime">{{ hours }}</div>
    </div>
  </div>
</template>

<style scoped>
.messageItem {
  display: flex;
  justify-content: flex-start;
  margin: 0.5rem 0;
  width: 100%;
}

.messageBox {
  background-color: #f0f7fc;
  border: 1px solid #046785;
  border-radius: 1rem;
  padding: 1rem;
  position: relative;
  box-shadow: 1px 3px 5px rgba(0, 0, 0, 0.1);
  max-width: 70%; /* Limita la larghezza massima */
  word-wrap: break-word; /* Spezza il testo lungo */
}

.messageHeader {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.optionsButton {
  background: none;
  border: none;
  font-size: 1.25rem;
  cursor: pointer;
  color: #046785;
  padding: 0;
  margin-left: 10px;
}

.optionsButton:hover {
  color: black;
}

.messageText {
  margin-top: 0.5rem;
  font-size: 1rem;
  font-family: "SF Pro", sans-serif;
  color: #04354c;
}

.messageTime {
  font-size: 0.75rem;
  color: #6c7a89;
  text-align: right;
  margin-top: 0.5rem;
}
</style>
