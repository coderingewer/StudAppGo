package main

import (
	"fmt"
	"net/http"
	"studapp/api/controllers"
	"studapp/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//User routers
	router.HandleFunc("/api/users/new", middlewares.SetMiddlewareJSON(controllers.CreateUser)).Methods("POST")
	router.HandleFunc("/api/users/getAll", middlewares.SetMiddlewareJSON(controllers.GetUsers)).Methods("GET")
	router.HandleFunc("/api/users/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetUser)).Methods("GET")
	router.HandleFunc("/api/users/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUser)).Methods("DELETE")
	router.HandleFunc("/api/users/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateUser))).Methods("POST")
	router.HandleFunc("/login", middlewares.SetMiddlewareJSON(controllers.Login)).Methods("POST")

	//Post routers
	router.HandleFunc("/api/posts/new", middlewares.SetMiddlewareJSON(controllers.CreatePost)).Methods("POST")
	router.HandleFunc("/api/posts/getAll", middlewares.SetMiddlewareJSON(controllers.GetPosts)).Methods("GET")
	router.HandleFunc("/api/posts/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetPost)).Methods("GET")
	router.HandleFunc("/api/posts/getByUserId/{userId}", middlewares.SetMiddlewareJSON(controllers.GetPostsByUserID)).Methods("GET")
	router.HandleFunc("/api/posts/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeletePost)).Methods("DELETE")
	router.HandleFunc("/api/posts/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePost))).Methods("POST")
	//Skill routers
	router.HandleFunc("/api/skills/new", middlewares.SetMiddlewareJSON(controllers.CreateSkill)).Methods("POST")
	router.HandleFunc("/api/skills/getAll", middlewares.SetMiddlewareJSON(controllers.GetSkills)).Methods("GET")
	router.HandleFunc("/api/skills/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetSkill)).Methods("GET")
	router.HandleFunc("/api/skills/getByUserId/{userId}", middlewares.SetMiddlewareJSON(controllers.GetSkillsByUserID)).Methods("GET")
	router.HandleFunc("/api/skills/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteSkill)).Methods("DELETE")
	router.HandleFunc("/api/skills/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateSkill))).Methods("POST")

	//City routers
	router.HandleFunc("/api/cities/new", middlewares.SetMiddlewareJSON(controllers.CreateCity)).Methods("POST")
	router.HandleFunc("/api/cities/getAll", middlewares.SetMiddlewareJSON(controllers.GetAllCities)).Methods("GET")

	//University routers
	router.HandleFunc("/api/universities/new", middlewares.SetMiddlewareJSON(controllers.CreateUniversity)).Methods("POST")
	router.HandleFunc("/api/universities/getByCityId/{cityId}", middlewares.SetMiddlewareJSON(controllers.GetByCityID)).Methods("GET")
	router.HandleFunc("/api/universities/getFacultyByUniId/{id}", middlewares.SetMiddlewareJSON(controllers.GetFacultyByUniID)).Methods("GET")
	router.HandleFunc("/api/universities/addFaculty/{id}/{facultyId}", middlewares.SetMiddlewareJSON(controllers.AddAFaculty)).Methods("POST")
	//Faculty routers
	router.HandleFunc("/api/faculties/new", middlewares.SetMiddlewareJSON(controllers.CreateFaculty)).Methods("POST")
	router.HandleFunc("/api/faculties/getAll", middlewares.SetMiddlewareJSON(controllers.GetFaculties)).Methods("GET")

	port := "8000"

	if port == "" {
		port = "8000"
	}
	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}

}
