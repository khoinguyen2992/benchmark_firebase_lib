let admin = require("firebase-admin");
let serviceAccount = require("./benchmark_firebase_credentials.json");

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: "https://benchmark-7f1d0.firebaseio.com"
});

let db = admin.database();

const maxKey = 30;
const maxRequest = 1000;

for (let index = 0; index < maxRequest; index++) {
  let ref = `nodejs/${index}`;
  let message = produceMessage('firebase-admin', index);

  console.log(`Set data at index ${index}`)

  db.ref(ref).set(message)
    .then(() => console.log(`Set at ${ref}`))
    .catch((err) => {
      console.log("Firebase-admin at ref: " + " Error: " + JSON.stringify(err));
    });
}

function produceMessage(identifier, count) {
  let object = {};

  for (let i = 0; i < maxKey; i++) {
    let key = `KEY_${identifier}_${count}_${i}`;
    let value = `VALUE_${count}_${i}`;
    object[key] = value;
  }

  return object;
}
