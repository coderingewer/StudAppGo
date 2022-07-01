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
	router.HandleFunc("/api/users/getByUni/{universtyId}", middlewares.SetMiddlewareJSON(controllers.GetUsersByUni)).Methods("GET")
	router.HandleFunc("/api/users/getByFaculty/{FacultyId}", middlewares.SetMiddlewareJSON(controllers.GetUsersByFaculty)).Methods("GET")
	router.HandleFunc("/api/users/getByFaculty/{departmentId}", middlewares.SetMiddlewareJSON(controllers.GetUsersByDepartmentID)).Methods("GET")

	router.HandleFunc("/api/posts/new", middlewares.SetMiddlewareJSON(controllers.CreatePost)).Methods("POST")
	router.HandleFunc("/api/posts/getAll", middlewares.SetMiddlewareJSON(controllers.GetPosts)).Methods("GET")
	router.HandleFunc("/api/posts/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetPost)).Methods("GET")
	router.HandleFunc("/api/posts/getByUserId/{userId}", middlewares.SetMiddlewareJSON(controllers.GetPostsByUserID)).Methods("GET")
	router.HandleFunc("/api/posts/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeletePost)).Methods("DELETE")
	router.HandleFunc("/api/posts/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePost))).Methods("POST")

	router.HandleFunc("/api/skills/new", middlewares.SetMiddlewareJSON(controllers.CreateSkill)).Methods("POST")
	router.HandleFunc("/api/skills/getAll", middlewares.SetMiddlewareJSON(controllers.GetSkills)).Methods("GET")
	router.HandleFunc("/api/skills/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetSkill)).Methods("GET")
	router.HandleFunc("/api/skills/getByUserId/{userId}", middlewares.SetMiddlewareJSON(controllers.GetSkillsByUserID)).Methods("GET")
	router.HandleFunc("/api/skills/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteSkill)).Methods("DELETE")
	router.HandleFunc("/api/skills/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateSkill))).Methods("POST")

	router.HandleFunc("/api/cities/new", middlewares.SetMiddlewareJSON(controllers.CreateCity)).Methods("POST")
	router.HandleFunc("/api/cities/getAll", middlewares.SetMiddlewareJSON(controllers.GetAllCities)).Methods("GET")

	router.HandleFunc("/api/universities/new", middlewares.SetMiddlewareJSON(controllers.CreateUniversity)).Methods("POST")
	router.HandleFunc("/api/universities/getByCityId/{cityId}", middlewares.SetMiddlewareJSON(controllers.GetByCityID)).Methods("GET")
	router.HandleFunc("/api/universities/addFaculty/{id}/{facultyId}", middlewares.SetMiddlewareJSON(controllers.AddAFaculty)).Methods("POST")
	router.HandleFunc("/api/universities/addDepartment/{id}/{facultyId}/{departmentId}", middlewares.SetMiddlewareJSON(controllers.AddADepartment)).Methods("POST")
	router.HandleFunc("/api/universities//getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetByID)).Methods("GET")
	router.HandleFunc("/api/universities/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUniversity)).Methods("DELETE")
	router.HandleFunc("/api/universities/deleteUniDepartment/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUniDepartmentByID)).Methods("DELETE")
	router.HandleFunc("/api/universities/deleteUniFaculty/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUniFacultyByID)).Methods("DELETE")

	router.HandleFunc("/api/faculties/new", middlewares.SetMiddlewareJSON(controllers.CreateFaculty)).Methods("POST")
	router.HandleFunc("/api/faculties/getAll", middlewares.SetMiddlewareJSON(controllers.GetFaculties)).Methods("GET")
	router.HandleFunc("/api/faculties/grtById/{id}", middlewares.SetMiddlewareJSON(controllers.GetFacultyByID)).Methods("GET")
	router.HandleFunc("/api/faculties/getByUniId/{universityId}", middlewares.SetMiddlewareJSON(controllers.GetFacultyByUniID)).Methods("GET")
	router.HandleFunc("/api/faculties/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteFacultyByID)).Methods("DELETE")

	router.HandleFunc("/api/departments/new", middlewares.SetMiddlewareJSON(controllers.CreateDepartment)).Methods("POST")
	router.HandleFunc("/api/depertments/getAll", middlewares.SetMiddlewareJSON(controllers.GetDepartments)).Methods("GET")
	router.HandleFunc("/api/depertments/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetDepartment)).Methods("GET")
	router.HandleFunc("/api/depertments/getByUniIdAndFacultyId/{universityId}/{facultyId}", middlewares.SetMiddlewareJSON(controllers.GetDepartmentByUniIDAndFacultyID)).Methods("GET")
	router.HandleFunc("/api/depertments/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteDepartment)).Methods("DELETE")

	router.HandleFunc("/api/images/upload", middlewares.SetMiddlewareJSON(controllers.ImgUpload)).Methods("POST")
	router.HandleFunc("/api/images/update/{imageId}", middlewares.SetMiddlewareJSON(controllers.UpdateImage)).Methods("POST")
	router.HandleFunc("/api/images/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteImageByUserID)).Methods("DELETE")

	router.HandleFunc("/api/amigos/new", middlewares.SetMiddlewareJSON(controllers.CreateAmigo)).Methods("POST")
	router.HandleFunc("/api/amigos/getByDESC", middlewares.SetMiddlewareJSON(controllers.GetAmigosByDESC)).Methods("GET")
	router.HandleFunc("/api/amigos/getByUserId", middlewares.SetMiddlewareJSON(controllers.GetAmigosByUserID)).Methods("GET")
	router.HandleFunc("/api/amigos/getByCityId/{cityId}", middlewares.SetMiddlewareJSON(controllers.GetAmigosByCityID)).Methods("GET")
	router.HandleFunc("/api/amigos/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteAmigo)).Methods("DELETE")

	port := "8000"

	if port != "8000" {
		port = "studappdemo.herokuapp.com"
	}
	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}

}
