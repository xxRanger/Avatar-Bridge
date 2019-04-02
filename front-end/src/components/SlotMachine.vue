<template>

    <section class="bg-light">
        <div class="row-container">
            <div class="role-info-container">
                <div class="avatar-info-container">
                    <div class="role-info-header">
                        Private chain avatar
                    </div>
                    <div v-if="avatarPrivateChain!==undefined">
                        <img class="avatar-img" :src="avatarPvcPath">
                    </div>
                </div>
                <div class="token-info-container">
                    <div class="role-info-header">
                        Private chain tokens
                    </div>
                    <div>
                        {{tokenPrivateChain}}
                    </div>
                </div>
            </div>
            <div class="game-container">
                <div class="game-instruction text-center">
                    <h3>How to play?</h3>
                    <br>
                    <br>
                    <b>click "run" button {{rewardRound}} times, </b><br>
                    <b>get {{gameReward}} token as a reward.</b>
                    <br>
                    <br>
                    <b>click "run" button {{avatarRound}} times</b><br>
                    <b>You can get an avatar or</b>
                    <br>
                    <b>upgrade your avatar</b>
                </div>
                <div class="game-content-container">
                    <div class='slot-machine'>
                        <div ref="slot1">
                            <div class='slot__wrap text-center'>
                                <div class='slot__item' v-bind:key="opt" v-for='opt in slots[0].items'>
                                    <img class="slot__icon" :src="opt">
                                </div>
                                <div class='slot__item slot__item--copy'>
                                    <img class="slot__icon" :src="slots[0].items[0]">
                                </div>
                            </div>
                        </div>
                        <div ref="slot2">
                            <div class='slot__wrap text-center' ref='slot2'>
                                <div class='slot__item' v-bind:key="opt" v-for='opt in slots[1].items'>
                                    <img class="slot__icon" :src="opt">
                                </div>
                                <div class='slot__item slot__item--copy'>
                                    <img class="slot__icon" :src="slots[1].items[0]">
                                </div>
                            </div>
                        </div>
                        <div ref="slot3">
                            <div class='slot__wrap text-center' ref='slot3'>
                                <div class='slot__item' v-bind:key="opt" v-for='opt in slots[2].items'>
                                    <img class="slot__icon" :src="opt">
                                </div>
                                <div class='slot__item slot__item--copy'>
                                    <img class="slot__icon" :src="slots[2].items[0]">
                                </div>
                            </div>
                        </div>
                    </div>


                    <div>
                        <form>
                            <div class="alpha-form">
                                <label class="form-label label slot-label" for="betInput">Wager: </label>
                                <input class="form-input" id="betInput" type="number" v-model="betAmount">
                                <button type="button" class="form-button alpha-button alpha-button-primary"
                                        @click="run">run
                                </button>
                            </div>
                        </form>
                    </div>
                    <div class="text-center">
                        <div class="inDouble text-center" v-if="doubled">
                            <b>Guaranteed winning Round !</b>
                        </div>
                        <span class="label">Slot reward:</span> <span v-if="reward!==0"
                                                                      class="slot-tag">{{reward}} </span>
                        <span
                                class="slotValue" v-else> {{reward}} </span>
                        <br>
                        <span class="label">Game cycle:</span> <span class="slotValue">{{gamePlayed}}</span>
                        <!--<br>-->
                        <!--<span class="label">Private chain tokens:</span> <span v-if="tokens!==0"-->
                        <!--class="slot-tag">{{tokenPrivateChain}} </span>-->
                        <!--<span class="slotValue" v-else> {{tokens}} </span>-->
                    </div>

                    <div class="consume text-center">
                        <div class="justify-content-md-center">
                            <div class="doubleHead">
                                <b> Want to get "three identical symbols"?</b>
                                <br>
                                <b>Pay {{doublePrice}} tokens </b>
                            </div>
                            <div class=" doubleButton text-center">
                                <button class="alpha-button alpha-button-primary" type="button" @click="consume">Have a
                                    try
                                </button>
                            </div>
                        </div>
                    </div>

                    <div class="consume text-center">
                        <div class="justify-content-md-center">
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>
<script>

    import bar from "@/assets/images/bar.png";
    import bell from "@/assets/images/bell.png";
    import cherry from "@/assets/images/cherry.png";
    import mango from "@/assets/images/mango.png";
    import orange from "@/assets/images/orange.png";
    import seven from "@/assets/images/seven.png";
    import {encodePrivateChainFunction, getPrivateChainBaseTxObject, signTx} from "../assets/js/tx";

    var next;
    export default {
        name: 'HelloWorld',
        props: {
            tokenPrivateChain: Number,
            tokenPublicChain: Number,
            avatarPrivateChain: Object,
            avatarPublicChain: Object,
            ws: WebSocket,
            renewDoubleRound: Number,
            avatarPvcPath: String,
            avatarPbcPath: String,
        },
        watch: {
            renewDoubleRound: function (newVal) {
                if (newVal !== 0) {
                    this.doubledRound += 2;
                    this.doubled = true;
                }
            },
            avatarPvcPath: function (v) {
                console.log(v);
            }
        },
        data() {
            return {
                avatarRound: 1,
                rewardRound: 1,
                gameReward: 100,
                doublePrice: 2,
                socketConstants: undefined,
                doubled: false,
                slots: [{
                    items: [
                        bar,
                        bell,
                        cherry,
                        mango,
                        orange,
                        seven
                    ],
                    itemsName: [
                        "bar",
                        "bell",
                        "cherry",
                        "mango",
                        "orange",
                        "seven"
                    ]
                }, {
                    items: [
                        cherry,
                        orange,
                        mango,
                        bar,
                        bell,
                        seven
                    ],
                    itemsName: [
                        "cherry",
                        "orange",
                        "mango",
                        "bar",
                        "bell",
                        "seven"
                    ]
                }, {
                    items: [
                        bell,
                        mango,
                        bar,
                        orange,
                        cherry,
                        seven
                    ],
                    itemsName: [
                        "bell",
                        "mango",
                        "bar",
                        "orange",
                        "cherry",
                        "seven"
                    ]
                }],
                opts: null,
                startedAt: null,
                betAmount: 0,
                result: undefined,
                reward: 0,
                gamePlayed: 0,
                doubledRound: 0,
                slotmachine: undefined,
                account: undefined,
            }
        },
        methods: {
            consume: function () {
                if (this.tokenPrivateChain < this.doublePrice) {
                    this.$store.state.notifyError('Tokens not enough');
                    return
                }
                this.$store.state.notifyInfo('Request has been received');
                let payload = {
                    gcuid: this.socketConstants.gcuid.consume,
                    type: this.socketConstants.consumeType.erc20Token,
                    amount:this.doublePrice,
                    Address: this.account.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            rewardToken: function () {
                this.$store.state.notifyInfo('Token on the way');
                let payload = {
                    gcuid: this.socketConstants.gcuid.reward,
                    address: this.account.address,
                    type: this.socketConstants.rewardType.erc20Token,
                    amount: this.gameReward,
                };
                this.ws.send(JSON.stringify(payload));
            },
            rewardAvatar: function () {
                this.$store.state.notifyInfo('Avatar on the way');
                let payload = {
                    gcuid: this.socketConstants.gcuid.reward,
                    address: this.account.address,
                    type: this.socketConstants.rewardType.avatar,
                };
                this.ws.send(JSON.stringify(payload));
            },
            upgradeAvatar: function () {  //TODO
                this.$store.state.notifyInfo('Upgrading Avatar');
                let payload = {
                    gcuid: this.socketConstants.gcuid.upgradeAvatar,
                    tokenId: this.avatarPrivateChain.tokenId,
                };
                this.ws.send(JSON.stringify(payload));
            },
            run: function () {
                if (this.opts) {
                    return
                }
                if (this.betAmount <= 0) {
                    this.$store.state.notifyError('Require wager > 0');
                    return
                }
                this.reward = 0;
                this.gamePlayed += 1;
                if (this.gamePlayed % this.rewardRound === 0) {
                    this.rewardToken()
                }

                if (this.gamePlayed % this.avatarRound === 0) {
                    console.log("avatar round");
                    if (this.avatarPublicChain === undefined && this.avatarPrivateChain === undefined) {
                        // no avatar
                        console.log("reward avatar");
                        this.rewardAvatar()
                    } else if (this.avatarPrivateChain !== undefined && this.avatarPrivateChain.avatarLevel < this.socketConstants.avatarMaxLevel) {
                        // avatar in private chain and has not reached max level
                        console.log("upgrade avatar");
                        this.upgradeAvatar()
                    }
                }

                let simulation = this.slotmachine.play(this.doubled);
                this.result = simulation.result;
                this.opts = this.slots.map((data, i) => {
                    let choice = this.getChoice(this.slots[i].itemsName, this.result[i]);
                    console.log(this.$refs);
                    let slot = this.$refs[`slot${i + 1}`];
                    //const choice = Math.floor( Math.random() * data.items.length )

                    return {
                        el: slot.querySelector('.slot__wrap'),
                        finalPos: choice * 180,
                        startOffset: 2000 + Math.random() * 500,
                        height: data.items.length * 180,
                        duration: 3000 + i * 700, // milliseconds
                        isFinished: false,
                    };

                });
                this.result = simulation.reward * this.betAmount;

                next(this.animate)

            },
            getChoice: function (arr, v) {
                for (let i = 0; i < arr.length; ++i) {
                    if (arr[i] === v) {
                        return i;
                    }
                }
                return 0;
            },
            animate: function (timestamp) {

                if (this.startedAt == null) {
                    this.startedAt = timestamp
                }

                const timeDiff = timestamp - this.startedAt;

                this.opts.forEach(opt => {

                    if (opt.isFinished) {
                        return
                    }

                    const timeRemaining = Math.max(opt.duration - timeDiff, 0);
                    const power = 3;
                    const offset = (Math.pow(timeRemaining, power) / Math.pow(opt.duration, power)) * opt.startOffset

                    // negative, such that slots move from top to bottom
                    const pos = -1 * Math.floor((offset + opt.finalPos) % opt.height)

                    opt.el.style.transform = "translateY(" + pos + "px)"

                    if (timeDiff > opt.duration) {
                        opt.isFinished = true
                    }

                });

                if (this.opts.every(o => o.isFinished)) {
                    this.reward = this.result;
                    if (this.doubled) {
                        this.doubledRound -= 1;
                        if (this.doubledRound === 0) {
                            this.doubled = false;
                        }
                    }
                    this.opts = null;
                    this.startedAt = null;
                    console.log('finished')
                } else {
                    next(this.animate)
                }

            },
        },
        created: function () {
            import('../scripts/slotmachine.js').then(s => {
                this.slotmachine = new s.Slotmachine();
            });
            this.socketConstants = this.$store.state.socketConstants;
            this.account = this.$store.state.account;
            // this.getTokens()
        },
        beforeCreate: function () {
            next = window.requestAnimationFrame ||
                window.webkitRequestAnimationFrame ||
                window.mozRequestAnimationFrame ||
                window.msRequestAnimationFrame ||
                window.oRequestAnimationFrame ||
                function (cb) {
                    window.setTimeout(cb, 1000 / 60)
                }
        },
    }

</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->