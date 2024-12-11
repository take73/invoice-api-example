package model

type ClientBankAccount struct {
    ID         int    // 銀行口座ID
    ClientID   int    // 紐づく取引先ID
    BankName   string // 銀行名
    BranchName string // 支店名
    AccountNumber string // 口座番号
    AccountName   string // 口座名
}