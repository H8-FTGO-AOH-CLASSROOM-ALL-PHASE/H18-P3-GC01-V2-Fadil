// seed products: _id, name, price, qty
db.products.insertMany([
    { _id: 1, name: "Laptop", price: 12000000, qty: 10 },
    { _id: 2, name: "Mouse", price: 150000, qty: 50 },
    { _id: 3, name: "Keyboard", price: 350000, qty: 30 },
    { _id: 4, name: "Monitor", price: 2500000, qty: 15 },
    { _id: 5, name: "Headset", price: 500000, qty: 25 }
])

// check seeded data
db.products.find()

//find product data

db.products.find()
