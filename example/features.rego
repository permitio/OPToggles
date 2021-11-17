package app.rbac

billing_users[users]{
  some user, i
  data.example.users[user].roles[i] == "billing"
  users := user
}

us_users[users]{
  some user
  data.example.users[user].location.country == "US"
  users := user
}