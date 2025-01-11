package model

type ElixirResp struct {
    Tag string
    Arabic string
    Translation string
}

type PerplexityResp struct {
    Translation string `json:"translation"`
    Examples []PerplexityRespExample `json:"examples"`
}

type PerplexityRespExample struct {
    Sentence string `json:"sentence"`
    Translation string `json:"translation"`
}

