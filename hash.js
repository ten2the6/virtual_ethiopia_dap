import { createHash, scryptSync, randomBytes, timeingSafeEqual } from 'node:crypto';

function hash(input) {
    return createHash('sha256').update(input).digest('hex')
}

let password = 'Tentothe6!'
const hash1 = hash(password)

console.log(hash1)


function signup(email, password) {
    const salt = randomBytes(16).toString('hex')
    const hashedPassword = scryptSync(password, salt, 64).toString('hex')

    const user = { email, password: `${salt}:${hashedPassword}`  }

}

function login(email, password) {
    const user = users.find(v => v.email === email);

    const [salt, key] = user.password.split(':');
    const hashedBuffer = scryptSync(password, salt, 64);

    const keyBuffer = Buffer.from(key, 'hex');
    const match = timeingSafeEqual(hashedBuffer, keyBuffer);

    if (match) {
        return 'login success!'
    } else {
        return 'login failed'
    }
}