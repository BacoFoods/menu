const fs = require('fs');

const input1 = fs.readFileSync('f1.txt', 'utf8').split('\n');
const input2 = fs.readFileSync('f2.txt', 'utf8').split('\n');

const queue = [];

for (let i = 0; i < input1.length; i++) {
    queue.push({i, val: input1[i]});
}

for (let i = 0; i < input2.length; i++) {
    pop(queue, input2[i]);
}


// valores en f1 y no en f2... valores en drive y no en excel
console.log(queue)

function pop(queue, val) {
    if (queue.length === 0) {
        return;
    }

    for (let i = 0; i < queue.length; i++) {
        if (`${queue[i].val}` === `${val}`) {
            queue.splice(i, 1);
            return;
        }
    }
}