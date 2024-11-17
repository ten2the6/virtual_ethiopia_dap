package api

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "virtual_ethiopia_dap/internal/blockchain"
)

type Server struct {
    chain  *blockchain.Chain
    router *mux.Router
}

// Response structure for all API responses
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// Request structures
type TransactionRequest struct {
    From   string  `json:"from"`
    To     string  `json:"to"`
    Amount float64 `json:"amount"`
}

type CitizenRegistrationRequest struct {
    Name        string `json:"name"`
    DateOfBirth string `json:"dateOfBirth"`
    PublicKey   string `json:"publicKey"`
}

type CitizenApprovalRequest struct {
    CitizenID   string `json:"citizenId"`
    ApproverKey string `json:"approverKey"`
}

type ElectionRequest struct {
    Name         string `json:"name"`
    DurationDays int    `json:"durationDays"`
}

type CandidateRequest struct {
    Name      string `json:"name"`
    PublicKey string `json:"publicKey"`
    Platform  string `json:"platform"`
}

type VoteRequest struct {
    CitizenPublicKey string `json:"citizenPublicKey"`
    CandidateID      string `json:"candidateId"`
}

func NewServer(chain *blockchain.Chain) *Server {
    server := &Server{
        chain:  chain,
        router: mux.NewRouter(),
    }
    server.setupRoutes()
    return server
}

func (s *Server) setupRoutes() {
    // Add middleware
    s.router.Use(loggingMiddleware)
    s.router.Use(corsMiddleware)

    // Blockchain endpoints
    s.router.HandleFunc("/blocks", s.handleGetBlocks).Methods("GET")
    s.router.HandleFunc("/transactions", s.handleAddTransaction).Methods("POST")
    
    // Citizen registry endpoints
    s.router.HandleFunc("/citizens/register", s.handleRegisterCitizen).Methods("POST")
    s.router.HandleFunc("/citizens/approve", s.handleApproveCitizen).Methods("POST")
    s.router.HandleFunc("/citizens", s.handleGetAllCitizens).Methods("GET")
    
    // Election endpoints
    s.router.HandleFunc("/elections/start", s.handleStartElection).Methods("POST")
    s.router.HandleFunc("/elections/candidates", s.handleRegisterCandidate).Methods("POST")
    s.router.HandleFunc("/elections/vote", s.handleCastVote).Methods("POST")
    s.router.HandleFunc("/elections/current", s.handleGetCurrentElection).Methods("GET")

    // Health check
    s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

func (s *Server) Start(port string) error {
    log.Printf("Starting API server on port %s\n", port)
    return http.ListenAndServe(":"+port, s.router)
}

// Handler implementations
func (s *Server) handleGetBlocks(w http.ResponseWriter, r *http.Request) {
    blocks := s.chain.GetBlocks()
    sendSuccess(w, blocks)
}

func (s *Server) handleAddTransaction(w http.ResponseWriter, r *http.Request) {
    var req TransactionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid transaction data", http.StatusBadRequest)
        return
    }

    tx, err := s.chain.AddTransaction(req.From, req.To, req.Amount)
    if err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendSuccess(w, tx)
}

func (s *Server) handleRegisterCitizen(w http.ResponseWriter, r *http.Request) {
    var req CitizenRegistrationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request data", http.StatusBadRequest)
        return
    }

    tx, err := s.chain.AddCitizenRegistration(req.Name, req.DateOfBirth, req.PublicKey)
    if err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendSuccess(w, tx)
}

func (s *Server) handleApproveCitizen(w http.ResponseWriter, r *http.Request) {
    var req CitizenApprovalRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request data", http.StatusBadRequest)
        return
    }

    tx, err := s.chain.ApproveCitizen(req.CitizenID, req.ApproverKey)
    if err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendSuccess(w, tx)
}

func (s *Server) handleStartElection(w http.ResponseWriter, r *http.Request) {
    var req ElectionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request data", http.StatusBadRequest)
        return
    }

    tx, err := s.chain.StartElection(req.Name, req.DurationDays)
    if err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendSuccess(w, tx)
}

func (s *Server) handleRegisterCandidate(w http.ResponseWriter, r *http.Request) {
    var req CandidateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request data", http.StatusBadRequest)
        return
    }

    tx, err := s.chain.RegisterCandidate(req.Name, req.PublicKey, req.Platform)
    if err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendSuccess(w, tx)
}

func (s *Server) handleCastVote(w http.ResponseWriter, r *http.Request) {
    var req VoteRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request data", http.StatusBadRequest)
        return
    }

    tx, err := s.chain.CastVote(req.CitizenPublicKey, req.CandidateID)
    if err != nil {
        sendError(w, err.Error(), http.StatusBadRequest)
        return
    }

    sendSuccess(w, tx)
}

func (s *Server) handleGetAllCitizens(w http.ResponseWriter, r *http.Request) {
    citizens := s.chain.GetAllCitizens()
    sendSuccess(w, citizens)
}

func (s *Server) handleGetCurrentElection(w http.ResponseWriter, r *http.Request) {
    election := s.chain.GetCurrentElection()
    if election == nil {
        sendError(w, "No active election", http.StatusNotFound)
        return
    }
    sendSuccess(w, election)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    sendSuccess(w, map[string]string{"status": "ok"})
}

// Helper functions
func sendError(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{
        Success: false,
        Error:   message,
    })
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{
        Success: true,
        Data:    data,
    })
}

// Middleware
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received %s request to %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}