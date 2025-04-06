package bo_v1

import (
	"encoding/json"
	"errors"
	models "lexicon/bo-api/beneficiary_ownership/v1/models"
	bo_v1_services "lexicon/bo-api/beneficiary_ownership/v1/services"
	"lexicon/bo-api/common/utils"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

func Router() *chi.Mux {

	r := chi.NewMux()
	r.Get("/search", searchHandler)
	r.Get("/detail/{id}", detailHandler)
	r.Get("/chart", chartHandler)
	r.Post("/chatbot", chatbotHandler)
	r.Post("/chatbot/references", chatbotReferenceHandler)
	return r
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	qp := r.URL.Query()
	query := qp.Get("query")
	rawSubjectType := qp.Get("subject_type")
	var subjectTypes []string

	if rawSubjectType != "" {
		subjectTypes = strings.Split(rawSubjectType, ",")
	}

	year := qp.Get("year")
	rawType := qp.Get("type")

	var caseTypes []string

	if rawType != "" {
		caseTypes = strings.Split(rawType, ",")
	}

	nations := strings.Split(qp.Get("nation"), ",")
	page := qp.Get("page")

	years := []string{}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("page must be a number"))
		return
	}

	if year != "" {
		yearsSplit := strings.Split(year, "-")
		if len(yearsSplit) != 2 {
			utils.WriteError(w, http.StatusBadRequest, errors.New("year must be in the format of year-year"))

			return
		}
		// genereate year between
		yearFrom, err := strconv.Atoi(yearsSplit[0])
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, errors.New("year must be in the format of year-year"))
			return
		}
		yearTo, err := strconv.Atoi(yearsSplit[1])

		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, errors.New("year must be in the format of year-year"))

			return
		}

		for i := yearFrom; i <= yearTo; i++ {
			years = append(years, strconv.Itoa(i))
		}

	}

	req := models.SearchRequest{
		Query:        query,
		SubjectTypes: subjectTypes,
		Years:        years,
		Types:        caseTypes,
		Nations:      nations,
		Page:         int64(pageInt),
	}

	response, err := bo_v1_services.Search(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if response.Data == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("data not found"))
		return
	}
	utils.WriteResponse(w, response, http.StatusOK)
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	response, err := bo_v1_services.GetDetail(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("data not found"))
		return
	}

	utils.WriteData(w, response, http.StatusOK)
}

func chartHandler(w http.ResponseWriter, r *http.Request) {
	response, err := bo_v1_services.GetChartData(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("data not found"))
		return
	}

	utils.WriteData(w, response, http.StatusOK)
}

func chatbotHandler(w http.ResponseWriter, r *http.Request) {
	req := models.ChatbotRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding request body")
		utils.WriteError(w, http.StatusBadRequest, errors.New("body is empty"))
		return
	}
	defer r.Body.Close()

	// Use ChatbotBaseURL from env variable
	chatbotBaseURL := os.Getenv("CHATBOT_BASE_URL")
	if chatbotBaseURL == "" {
		chatbotBaseURL = ""
	}

	chatReq, err := http.NewRequestWithContext(r.Context(), "POST", chatbotBaseURL+"/chatbot/user_message", nil)
	if err != nil {
		log.Error().Err(err).Msg("Error creating request")
		utils.WriteError(w, http.StatusNotFound, errors.New("data not found"))
		return
	}

	chatReq.Header.Set("Accept", "application/json")
	chatReq.Header.Set("Content-Type", "application/json")

	params := url.Values{}
	params.Add("thread_id", req.ThreadID)
	params.Add("user_message", req.UserMessage)

	chatReq.URL.RawQuery = params.Encode()

	log.Info().Msgf("Chatbot request: %s", chatReq.URL.String())
	chatResp, err := utils.Client.Do(chatReq)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("data not found"))
		return
	}

	defer chatResp.Body.Close()

	response := models.ChatbotResponse{}

	err = json.NewDecoder(chatResp.Body).Decode(&response)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding response body")
		utils.WriteError(w, http.StatusBadRequest, errors.New("body is empty"))
		return
	}

	utils.WriteData(w, response, http.StatusOK)
}

func chatbotReferenceHandler(w http.ResponseWriter, r *http.Request) {
	req := models.ChatbotReferenceRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding request body")
		utils.WriteError(w, http.StatusBadRequest, errors.New("body is empty"))
		return
	}
	urls, err := bo_v1_services.GetUrlByCaseNumber(r.Context(), req.CaseNumbers)
	if err != nil {
		log.Error().Err(err).Msg("Error getting urls")
		utils.WriteError(w, http.StatusNotFound, errors.New("data not found"))
		return
	}
	utils.WriteData(w, urls, http.StatusOK)
}
