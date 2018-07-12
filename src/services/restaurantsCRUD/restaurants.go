package restaurantsCRUD

import (
	"github.com/Nastya-Kruglikova/cool_tasks/src/models"
	"github.com/Nastya-Kruglikova/cool_tasks/src/services/common"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

type successCreate struct {
	Status string            `json:"message"`
	Result models.Restaurant `json:"result"`
}

type successDelete struct {
	Status string `json:"message"`
}

func Get(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if val, ok := query["id"]; ok {
		id, err := uuid.FromString(val[0])
		if err != nil {
			common.SendNotFound(w, r, "ERROR: Invalid ID", err)
			return
		}
		items, err := models.GetRestaurantsByID(id)

		if err != nil {
			common.SendNotFound(w, r, "ERROR: Can't get items", err)
			return
		}

		common.RenderJSON(w, r, items)
	}

	items, err := models.GetRestaurantsByQuery(query)

	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get items", err)
		return
	}

	common.RenderJSON(w, r, items)
}

func GetFromTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	tripID, err := uuid.FromString(params["id"])
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong trip ID (can't convert string to uuid)", err)
		return
	}

	restaurants, err := models.GetRestaurantsFromTrip(tripID)
	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get restaurants by trip ID", err)
		return
	}

	common.RenderJSON(w, r, restaurants)
}

func Post(w http.ResponseWriter, r *http.Request) {

	var (
		newItem    models.Restaurant
		resultItem models.Restaurant
	)

	err := r.ParseForm()

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Can't parse POST Body", err)
		return
	}

	newItem.Name = r.Form.Get("name")
	newItem.Location = r.Form.Get("location")
	newItem.Description = r.Form.Get("description")
	newItem.Prices, err = strconv.Atoi(r.Form.Get("prices"))
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Invalid prices field", err)
		return
	}
	newItem.Stars, err = strconv.Atoi(r.Form.Get("stars"))
	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Invalid stars field", err)
		return
	}

	resultItem, err = models.AddRestaurant(newItem)

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Can't add new item", err)
		return
	}

	common.RenderJSON(w, r, successCreate{Status: "201 Created", Result: resultItem})
}

func Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	itemID, err := uuid.FromString(params["id"])

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong item ID (can't convert string to uuid)", err)
		return
	}

	err = models.DeleteRestaurantsFromDB(itemID)

	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't delete this item", err)
		return
	}

	common.RenderJSON(w, r, successDelete{Status: "204 No Content"})
}