<script>
import api from '../api'

export default {
    name: 'index',

    props: {
        msg: String
    },

    data() {
        return {
            error: false,
            list: [],
            kanji: null,
            radicalFrame: 0,
        }
    },

    created() {
        this.getGrade(1);
    },

    methods: {

        _handleError: function(error) {
            this.error = error;
        },

        getGrade: function(grade) {
            api.get('/search/advanced?grade='+grade).then( (response) => {
                this.list = response;
                this.getKanji(this.list[Math.floor(Math.random()*this.list.length)].ka_utf);
            }).catch( this._handleError );
        },

        getKanji: function(utf) {
            api.get('/kanji/'+utf).then( (response) => {
                console.log(response);
                this.kanji = response;
            }).catch( this._handleError );
        },

        play: function(event) {
            let parent = event.target.parentElement;
            var audio = parent.children[1];
            audio.play();
        },
        startRadicalAnimation: function() {
            setTimeout(() => {
                this.radicalFrame = 1;
            }, 1250, false);
            setTimeout(() => {
                this.radicalFrame = 2;
            }, 2500, false);
            setTimeout(() => {
                this.radicalFrame = 0;
            }, 3750, false);
        },
    },
}
</script>

<template>
    <div class="container">
        <div v-if="error" class="alert alert-danger">
            {{ error }}
        </div>
        <div v-if="kanji" class="row">
            <div class="col-sm mt-4 mb-2">
                <img width="230" :src="kanji.gothic_source">
            </div>
            <div class="col-sm mt-4">
                    <h1 class="display-4">{{ kanji.meaning }}</h1>
                    <span class="badge float-right">kunyomi</span>
                    <hr>
                    <p class="lead">{{ kanji.kunyomi_ja }}<br>{{ kanji.kunyomi }}</p>
                    <span class="badge float-right">onyomi</span>
                    <hr>
                    <p class="lead">{{ kanji.onyomi_ja }}<br v-if="kanji.onyomi_ja">{{ kanji.onyomi }}</p>
            </div>
        </div>
        <div v-if="kanji" class="row">
            <div class="col-sm">
                <div class="row">
                    <div class="col">
                    <div class="card radical">
                        <div class="card-header">
                            <h5 class="card-title">Radical</h5>
                        </div>
                        <div class="card-body">
                            <p class="card-text">
                                <img class="radical-character" :src="kanji.rad_char_source">
                                <span class="radical-meaning">{{ kanji.rad_meaning }}</span>
                                <div class="rad-anim float-right" @click="startRadicalAnimation">
                                    <img :class="{'opaque': radicalFrame == 0}"  :src="kanji.rad_anim_frame_0">
                                    <img :class="{'opaque': radicalFrame == 1}"  :src="kanji.rad_anim_frame_1">
                                    <img :class="{'opaque': radicalFrame == 2}"  :src="kanji.rad_anim_frame_2">
                                </div>
                            </p>
                        </div>
                    </div>
                    </div>
                    <div class="col">
                    <div class="card">
                        <div class="card-header">
                            <h5 class="card-title">Hint</h5>
                        </div>
                        <div class="card-body">
                            <p class="card-text hint" v-html="kanji.hint"></p>
                        </div>
                    </div>
                    </div>
                </div>
            </div>
            <div class="col-sm">
                <div class="card">
                    <div class="card-header">
                        <h5 class="card-title">Examples</h5>
                    </div>
                    <ul class="list-group list-group-flush">
                        <li class="list-group-item" v-for="ex in kanji.examples">
                            {{ ex.english }} / {{ ex.japanese }}
                            <span><i class="fa fa-play pull-right" @click="play"></i>
                                <audio>
                                    <source v-if="ex.opus" type="audio/opus" preload="auto" :src="ex.opus">
                                    <source v-if="ex.aac" type="audio/aac" preload="auto" :src="ex.aac">
                                    <source v-if="ex.ogg" type="audio/ogg" preload="auto" :src="ex.ogg">
                                    <source v-if="ex.mp3" type="audio/mp3" preload="auto" :src="ex.mp3">
                                </audio>
                            </span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">

.rad-anim img {
    cursor: pointer;
    width:64px; 
    height:66px; 
    -webkit-transition: opacity 1s ease-in-out;
    -moz-transition: opacity 1s ease-in-out;
    -o-transition: opacity 1s ease-in-out;
    transition: opacity 1s ease-in-out;
    opacity:0;
    display:none;
    -ms-filter:"progid:DXImageTransform.Microsoft.Alpha(Opacity=0)";
    filter: alpha(opacity=0);
}

.rad-anim img.opaque {
    opacity:1;
    display:block;
    -ms-filter:"progid:DXImageTransform.Microsoft.Alpha(Opacity=100)";
    filter: alpha(opacity=1);
}

.radical img {
    max-width:50px;
    max-height:50px;
}

.hint img {
    max-width:12px;
    max-height:12px;
}

.radical-character {
    margin-left:0px;
    width:31px;
    height:32px;
    margin-left:2px;
}

div.radical-meaning{
    white-space:normal;
    width:143px;
}
</style>
