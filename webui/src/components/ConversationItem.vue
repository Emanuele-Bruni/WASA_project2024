<script>
export default {
    props : {
        photo: String,
        title: String,
        isGroup: Boolean,
        lastMessage: {
          type: Object,
          default: null
        }
    },
}
</script>


<template>
  <div class="conversationCard">
    <img :src="'data:image/png;base64,' + photo" class="avatar" /> 


    <div class="info">
      <h3>{{ title }}</h3>
      <div v-if="lastMessage" class="lastMessage">
        <span class="hours">{{ lastMessage.hours }} | </span>
        <span class="author">{{ lastMessage.author }}:</span>
        <span class="text">{{ lastMessage.text }}</span>
      </div>
    </div>

    <span class="statusDot" :class="isGroup ? 'group' : 'private'"></span>
  </div>
</template>


<style scoped>
.lastMessage {
  font-size: 14px;
  color: #a0a0a0; /* grigio chiaro */
  font-style: italic;
  display: flex;
  flex-wrap: nowrap;
  gap: 5px;
  margin-top: 2px;
}

.lastMessage .author {
  font-weight: bold;
}

.lastMessage .hours {
  margin-left: auto;
  font-style: italic;
  align-items: right;
}
.conversationCard {
  display: flex;
  align-items: center;
  position: relative;
  padding: 0.75rem;
  margin: 0.5rem 0;
  background-color: #f9f9f9;
  border-radius: 12px;
  border: 0.5px solid rgb(4, 103, 133);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}
.conversationCard:hover{
    border: 2px solid rgb(4,103,133);
    box-shadow: 4px 3px 3px rgb(4, 103, 133);
}
.conversationCard:active{
  border: 2px solid rgb(4,103,133);
  box-shadow: 
	inset 0px -2px 12px -3px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 2px 12px -3px rgba(11, 51, 68, 0.5),
	inset 2px 0px 12px -3px rgba(11, 51, 68, 0.5),
	inset -2px 0px 12px -3px rgba(11, 51, 68, 0.5);
}
.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 1rem;
}

.info h3 {
  margin: 0;
  font-family: "SF Pro";
  color: rgb(4, 103, 133);
  font-size: 30px;
}

/* Pallino in basso a destra */
.statusDot {
  position: absolute;
  bottom: 10px;
  right: 10px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.statusDot.group {
  background-color: red;
}

.statusDot.private {
  background-color: gray;
}
</style>
