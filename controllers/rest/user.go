package rest

import (
	"project/controllers"
	"project/models"
)

type UserController struct {
	controllers.Controller
}


func (c *UserController) Index(page int, pagesize int) {
    
    // if c.Session == nil {
    //     c.Result["code"] = "auth error"
    //     return
    // }
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    var args []interface{}
    
    _name := c.Query("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
    }
    _passwd := c.Query("passwd")
    if _passwd != "" {
        args = append(args, models.Where{Column:"passwd", Value:_passwd, Compare:"like"})
    }
    
    _email := c.Query("email")
    if _email != "" {
        args = append(args, models.Where{Column:"email", Value:_email, Compare:"like"})
    }
    _startdate := c.Query("startdate")
    _enddate := c.Query("enddate")
    if _startdate != "" && _enddate != "" {        
        var v [2]string
        v[0] = _startdate
        v[1] = _enddate  
        args = append(args, models.Where{Column:"date", Value:v, Compare:"between"})    
    } else if  _startdate != "" {          
        args = append(args, models.Where{Column:"date", Value:_startdate, Compare:">="})
    } else if  _enddate != "" {          
        args = append(args, models.Where{Column:"date", Value:_enddate, Compare:"<="})            
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


func (c *UserController) Read(id int64) {
    
    // if c.Session == nil {
    //     c.Result["code"] = "auth error"
    //     return
    // }
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	item := manager.Get(id)

    c.Set("item", item)
}


func (c *UserController) Insert(item *models.User) {
    
	conn := c.NewConnection()
    
	manager := models.NewUserManager(conn)
	manager.Insert(item)

    id := manager.GetIdentity()
    c.Result["id"] = id
    item.Id = id
}

func (c *UserController) Update(item *models.User) {
    
    // if c.Session == nil {
    //     c.Result["code"] = "auth error"
    //     return
    // }
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Update(item)
}

func (c *UserController) Delete(item *models.User) {
    
    // if c.Session == nil {
    //     c.Result["code"] = "auth error"
    //     return
    // }
    
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Delete(item.Id)
}
