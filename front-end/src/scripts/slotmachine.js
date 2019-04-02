var config, pattern;

function _play(cheat) {
    var icon1 = reel(random(), pattern.reel1Segment);
    var icon2 = reel(random(), pattern.reel2Segment);
    var icon3 = reel(random(), pattern.reel3Segment);
    if (cheat) {
        var reward = judge(icon1, icon1, icon1)
        return {
            "reward": reward,
            "result": [icon1, icon1, icon1]
        }
    } else {
        var reward = judge(icon1, icon2, icon3);
        return {
            "reward": reward,
            "result": [icon1, icon2, icon3]
        }
    }
}

function judge(icon1, icon2, icon3) {
    var query = icon1 + icon2 + icon3;
    for (var key in config.award) {
        let match = query.match(key);
        if (match && match.index == 0) {
            return config.award[key];
        }
    }
    return 0;
}

function switchPattern(difficulty) {
    for (var i = 0; i < config.difficultySegment.length; ++i) {
        if (difficulty <= config.difficultySegment[i]) {
            return config.pattern[i];
        }
    }
    return config.pattern[config.difficultySegment.length - 1];
}

function reel(num, reelSegmeent) {
    for (var i = 0; i < reelSegmeent.length; ++i) {
        if (num <= reelSegmeent[i]) {
            return config.icon[i];
        }
    }
    return config.icon[reelSegmeent.length - 1];
}


function random() {
    return Math.floor(Math.random() * 20) + 1; // 1~20;
}

function Slotmachine() {

    config = require('../etc/slotmachineConfig.json');
    pattern = config.pattern3;


    this.difficulty = config.maxDifficulty;

    this.reduceDifficulty = function() {
        console.log("difficulty before:", this.difficulty);
        this.difficulty -= 1;
        pattern = switchPattern(this.difficulty);
        console.log("difficulty after: ", this.difficulty)
        return this.difficulty;
    }

    this.play = function(cheat) {
        try {
            return _play(cheat);
        } catch (err) {
            console.log(err.message);
        }
    }
}

export {
    Slotmachine
}