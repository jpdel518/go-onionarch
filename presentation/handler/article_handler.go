package handler

import (
	"encoding/json"
	"github.com/jpdel518/go-onionarch/domain/model"
	"github.com/jpdel518/go-onionarch/usecase"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	u usecase.ArticleUsecase
}

func NewArticleHandler(usecase usecase.ArticleUsecase) *Handler {
	return &Handler{usecase}
}

// func (h *Handler) RestHandle(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		h.fetch(w, r)
// 	} else if r.Method == http.MethodPost {
// 		h.store(w, r)
// 	} else if r.Method == http.MethodPut {
// 		h.update(w, r)
// 	} else if r.Method == http.MethodDelete {
// 		h.delete(w, r)
// 	}
// }

// ShowById get an article by id
func (h *Handler) ShowById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		// parameters
		id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get data from repository via usecase
		ctx := r.Context()
		article, err := h.u.GetByID(ctx, id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// conversion to json
		data, err := json.Marshal(article)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		if _, err := w.Write(data); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ShowByTitle get an article by title
func (h *Handler) ShowByTitle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		// parameters
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		title := r.PostFormValue("title")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get data from repository via usecase
		ctx := r.Context()
		article, err := h.u.GetByTitle(ctx, title)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// conversion to json
		data, err := json.Marshal(article)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		if _, err := w.Write(data); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Fetch get articles
func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	// get
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		// parameters
		numQuery := r.URL.Query().Get("num")
		var num int64 = 10
		if len(numQuery) >= 0 {
			num, _ = strconv.ParseInt(numQuery, 10, 64)
		}

		// get data from repository via usecase
		ctx := r.Context()
		articles, err := h.u.Fetch(ctx, num)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// conversion to json
		data, err := json.Marshal(articles)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		if _, err := w.Write(data); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Update update an article
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		w.Header().Set("Content-Type", "application/json")

		// parameters
		// err := r.ParseForm()
		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		ar := model.Article{}
		err := json.NewDecoder(r.Body).Decode(&ar)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// update article
		ctx := r.Context()
		err = h.u.Update(ctx, &ar)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte("successfully updated")); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Store register an article
func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		// parameters
		// err := r.ParseForm()
		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		ar := model.Article{}
		err := json.NewDecoder(r.Body).Decode(&ar)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// store article
		ctx := r.Context()
		err = h.u.Store(ctx, &ar)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte("successfully registered")); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Delete delete an article by id
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")

		// parameters
		id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete article
		ctx := r.Context()
		err = h.u.Delete(ctx, id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte("successfully deleted")); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
