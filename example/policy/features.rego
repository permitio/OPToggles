package app.rbac

billing_users[users]{
  some user, i
  data.users[user].roles[i] == "billing"
  users := user
}

us_users[users]{
  some user
  data.users[user].location.country == "US"
  users := user
}