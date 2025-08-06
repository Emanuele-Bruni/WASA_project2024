<script>
export default {
    props : {
        photo: String,
        title: String,
        isGroup: Boolean,
        members: Array,
    },
    data() {
        return {
            editableTitle: this.title
        };
    }
}
</script>
<template>
    <div class="modal-overlay">
        <div class="modal-content">
            <!-- Foto -->
             <div class="avatar-wrapper">
                <img :src="`data:image/png;base64,${photo.replace(/\s/g, '')}`" class="avatar" />
                <!-- Modificabile solo se è un gruppo -->
                 <button v-if="isGroup" class="edit-button"> Photo </button>
             </div>

            <!-- Titolo -->
              <h2>
                <span v-if="isGroup" class="edit-title-span">
                  <form id="SetGroupName" @submit.prevent="$emit('groupNameChanged', editableTitle)">
                    <input v-model="editableTitle" :placeholder="title" />
                    <button type="submit"> Modifica </button>
                  </form>

                </span>
                <span v-else class="singleTitle"> {{ title }}</span>
              </h2>

            <!-- Lista Membri -->
              <h3> Membri :</h3>
              <ul class="members-list">
                    <li v-for="member in members" :key="member.id"> {{ member.name }}</li>
              </ul>

            <!-- Pulsanti se è un gruppo-->
             <div v-if="isGroup" class="actions">
                <button> Aggiungi Membro </button>
                <button type="submit" @click="$emit('left')"> Esci dal Gruppo </button>
             </div>

             <button class="close-button" @click="$emit('close')"> Chiudi </button>
        </div>
    </div>
</template>

<style scoped>
.singleTitle{
  color: #046785;
}
h3{
  color: rgb(4, 103, 133);
  font-family: "SF Pro";
}
.edit-title-span{
  display: flex;
  flex-direction: column;
}
.edit-title-span button{
  align-items: center;
  width: 100%;
	border: none;
	border-radius: 8px;
	font-family: "SF Pro";
  background-color: rgb(4, 103, 133);
	color: white;
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  width: 450px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}
.avatar-wrapper {
  position: relative;
  width: 100px;
  margin: 0 auto;
}
.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
}
.edit-button {
  position: absolute;
  bottom: 0;
  right: 0;
  background: #046785;
  color: white;
  border: none;
  cursor: pointer;
}
.actions {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.actions button{
  padding: 5px 10px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-family: "SF Pro";
  background-color: rgb(4, 103, 133);
	color: white;
}
.close-button {
	text-align: center;
	font-family: "SF Pro";
	color: rgb(243, 240, 240); 
	background-color: rgb(235, 85, 85);
	border: 2px solid rgb(243, 32, 32);
  border-radius: 8px;
	padding: 0.2rem;
	width: 100%;
	margin-top: 0.3rem;
	box-shadow: 1px 3px 0 rgb(243, 32, 32);
}
.close-button:hover{
	background-color:rgb(243, 32, 32);
}
.close-button:active{
  box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.members-list {
  list-style: none;
  padding: 0;
  color: rgb(4, 66, 85); 
}
</style>