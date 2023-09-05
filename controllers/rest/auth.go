package rest

import (
	"project/controllers"
	"project/models"
)

type AuthController struct {
	controllers.Controller
}

func (c *AuthController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewAuthManager(conn)

    var args []interface{}

    name := c.Query("user")
    if name != "" {
        args = append(args, models.Where{Column:"user", Value:name, Compare:"="})
    }

    passwd := c.Query("token")
    if passwd != "" {
        args = append(args, models.Where{Column:"token", Value:passwd, Compare:"like"})
    }

    // email := c.Query("email")
    // if email != "" {
    //     args = append(args, models.Where{Column:"email", Value:email, Compare:"like"})
    // }
    startdate := c.Query("startdate")
    enddate := c.Query("enddate")
    if startdate != "" && enddate != "" {
        var v [2]string
        v[0] = startdate
        v[1] = enddate
        args = append(args, models.Where{Column:"date", Value:v, Compare:"between"})
    } else if  startdate != "" {
        args = append(args, models.Where{Column:"date", Value:startdate, Compare:">="})
    } else if  enddate != "" {
        args = append(args, models.Where{Column:"date", Value:enddate, Compare:"<="})
    }
    
    if page != 0 && pagesize != 0 {
        args = append(args, models.Paging(page, pagesize))
    }

    orderby := c.Query("orderby")
    if orderby == "desc" {
        // if page != 0 && pagesize != 0 {
            orderby = "id desc"
        // }
    } else {
		orderby = ""
	}

    if orderby != "" {
        args = append(args, models.Ordering(orderby))
    }

	items := manager.Find(args)
	c.Set("items", items)

    total := manager.Count(args)
	c.Set("total", total)
}

func (c *AuthController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewAuthManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}

func (c *AuthController) Insert(item *models.Auth) {
	conn := c.NewConnection()

	manager := models.NewAuthManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *AuthController) Update(item *models.Auth) {
	conn := c.NewConnection()

	manager := models.NewAuthManager(conn)
	manager.Update(item)
}

func (c *AuthController) Delete(item *models.Auth) {
	conn := c.NewConnection()

	manager := models.NewAuthManager(conn)
	manager.Delete(item.Id)
}