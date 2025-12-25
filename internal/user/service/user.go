package service

import (
	"cashapp/core"
	"cashapp/internal/user/models"
	"cashapp/internal/user/repository"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	repository repository.Repo
	config     *core.Config
}

func New(r repository.Repo, c *core.Config) *UserService {
	return &UserService{
		repository: r,
		config:     c,
	}
}

func (s *UserService) CreateUser(req core.CreateUserRequest) core.Response {
	user, err := s.repository.Users.FindByTag(req.Tag)

	if err == nil {
		return core.Error(errors.New("cash tag taken"), core.String("cash tag has already been taken"))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return core.Error(err, nil)
	}

	user = &models.User{
		Tag:       req.Tag,
		KYCLevel:  0,
		KYCStatus: models.KYCStatusPending,
		RiskScore: 0,
	}

	if err := s.repository.Users.Create(user); err != nil {
		return core.Error(err, nil)
	}

	wallet, err := s.repository.Wallets.Create(user.ID)
	if err != nil {
		return core.Error(err, nil)
	}

	user.Wallets = append(user.Wallets, *wallet)
	return core.Success(&map[string]interface{}{
		"user": user,
	}, core.String("user created successfully"))
}

func (s *UserService) GetUser(tag string) core.Response {
	user, err := s.repository.Users.FindByTag(tag)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Error(err, core.String("user not found"))
		}
		return core.Error(err, nil)
	}

	// Fetch primary wallet
	wallet, err := s.repository.Wallets.FindPrimaryWallet(user.ID)
	if err == nil {
		user.Wallets = append(user.Wallets, *wallet)
	}

	return core.Success(&map[string]interface{}{
		"user": user,
	}, nil)
}

func (s *UserService) InitVerification(req core.VerifyIdentityRequest) core.Response {
	user, err := s.repository.Users.FindByID(req.UserID)
	if err != nil {
		return core.Error(err, core.String("user not found"))
	}

	// Store document metadata
	doc := &models.IdentityDocument{
		UserID: user.ID,
		Type:   req.DocumentType,
		URL:    req.DocumentURL,
		Status: "pending",
	}

	if err := s.repository.IdentityDocuments.Create(doc); err != nil {
		return core.Error(err, core.String("failed to create document record"))
	}

	// In real world, call Stripe Identity / Onfido API here
	// For MVP, we simulate a "pending" state and return a fake session URL
	fmt.Printf("Mock: Initializing verification for user %d with doc %s\n", user.ID, req.DocumentType)

	return core.Success(&map[string]interface{}{
		"session_url": fmt.Sprintf("https://verify.mockprovider.com?user=%d", user.ID),
		"status":      "pending_verification",
		"document_id": doc.ID,
	}, nil)
}

func (s *UserService) HandleIdentityWebhook(req core.IdentityWebhookRequest) core.Response {
	user, err := s.repository.Users.FindByID(req.UserID)
	if err != nil {
		return core.Error(err, core.String("user not found"))
	}

	// If DocumentID is provided, update the specific document
	if req.DocumentID != 0 {
		doc, err := s.repository.IdentityDocuments.FindByID(req.DocumentID)
		if err == nil {
			if strings.EqualFold(req.Status, "passed") {
				doc.Status = "verified"
			} else {
				doc.Status = "rejected"
			}
			s.repository.IdentityDocuments.Update(doc)
		}
	}

	if strings.EqualFold(req.Status, "passed") {
		// Mock Sanctions Screening (AML)
		// In a real system, we'd check against OFAC/SDN lists here.
		// For simulation, let's say if the user tag contains "sanction", we flag them.
		if strings.Contains(strings.ToLower(user.Tag), "sanction") {
			user.KYCStatus = models.KYCStatusRejected
			user.RiskScore = 100 // Critical risk
			core.Log.Warn("User flagged during Sanctions Screening", zap.String("tag", user.Tag))
		} else {
			user.KYCStatus = models.KYCStatusVerified
			user.KYCLevel = 2   // Full verified
			user.RiskScore = 10 // Low risk
		}
	} else {
		user.KYCStatus = models.KYCStatusRejected
		user.RiskScore = 90 // High risk
	}

	if err := s.repository.Users.Update(user); err != nil {
		return core.Error(err, nil)
	}

	return core.Success(&map[string]interface{}{
		"user_id":    user.ID,
		"kyc_status": user.KYCStatus,
	}, nil)
}

func (s *UserService) LinkFundingSource(req core.LinkFundingSourceRequest) core.Response {
	// Mock: Retrieve payment method details from Stripe using PaymentMethodID
	// stripe.PaymentMethod.Get(req.PaymentMethodID)

	fs := &models.FundingSource{
		UserID:     req.UserID,
		Type:       req.Type,
		ProviderID: req.PaymentMethodID,
		Last4:      "4242", // Mock
		Brand:      "Visa", // Mock
	}

	if err := s.repository.FundingSources.Create(fs); err != nil {
		return core.Error(err, core.String("failed to link funding source"))
	}

	return core.Success(&map[string]interface{}{
		"funding_source": fs,
	}, core.String("funding source linked"))
}

func (s *UserService) Deposit(req core.DepositRequest) core.Response {
	// 1. Validate Funding Source
	fs, err := s.repository.FundingSources.FindByID(req.FundingSourceID)
	if err != nil {
		return core.Error(err, core.String("funding source not found"))
	}
	if fs.UserID != req.UserID {
		return core.Error(errors.New("unauthorized"), core.String("funding source does not belong to user"))
	}

	// 2. Mock Stripe Charge (Synchronous for now)
	// stripe.PaymentIntents.Create(...)
	fmt.Printf("Mock: Charging %d cents from %s for User %d\n", req.Amount, fs.ProviderID, req.UserID)
	// Simulate success

	// 3. Credit Wallet
	wallet, err := s.repository.Wallets.FindPrimaryWallet(req.UserID)
	if err != nil {
		return core.Error(err, core.String("wallet not found"))
	}

	wallet.Balance += req.Amount
	if err := s.repository.Wallets.Update(wallet); err != nil {
		// In real world, we need to reverse the charge here!
		return core.Error(err, core.String("failed to update wallet balance"))
	}

	return core.Success(&map[string]interface{}{
		"new_balance": wallet.Balance,
		"currency":    wallet.Currency,
	}, core.String("deposit successful"))
}

func (s *UserService) AddFriend(req core.CreateFriendshipRequest) core.Response {
	// Check if already friends
	if _, err := s.repository.Friendships.Find(req.UserID, req.FriendID); err == nil {
		return core.Error(nil, core.String("already friends"))
	}

	// Create friendship (bi-directional for simplicity, or two rows)
	// For MVP, one row implies check both ways or create 2 rows.
	// Let's create two rows to simplify querying.

	f1 := models.Friendship{
		UserID:   req.UserID,
		FriendID: req.FriendID,
		Status:   "accepted", // Auto-accept for MVP
	}
	f2 := models.Friendship{
		UserID:   req.FriendID,
		FriendID: req.UserID,
		Status:   "accepted",
	}

	if err := s.repository.Friendships.Create(&f1); err != nil {
		return core.Error(err, core.String("failed to add friend"))
	}
	s.repository.Friendships.Create(&f2) // Ignore error for 2nd row for now

	return core.Success(nil, core.String("friend added"))
}
