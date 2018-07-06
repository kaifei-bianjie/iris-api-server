package irishub

import (
	irisProtoc "github.com/irisnet/irishub-rpc/codegen/server/model"
	"github.com/irisnet/irishub-server/rpc"
	"github.com/irisnet/irishub-server/rpc/vo"
	"golang.org/x/net/context"
)

type CandidateListHandler struct {

}

func (c CandidateListHandler) Handler(ctx context.Context, req *irisProtoc.CandidateListRequest) (
	[]*irisProtoc.Candidate, error) {
	
	reqVO := c.BuildRequest(req)
	resVO, err := candidateService.List(reqVO)
	
	if err.IsNotNull() {
		return nil, rpc.ConvertIrisErrToGRPCErr(err)
	}
	
	return c.BuildResponse(resVO), nil
}

func (c CandidateListHandler) BuildRequest(req *irisProtoc.CandidateListRequest) vo.CandidateListReqVO {
	reqVO := vo.CandidateListReqVO{
		Address: req.GetAddress(),
		Page: req.GetPage(),
		PerPage: req.GetPerPage(),
		Sort: req.GetSort(),
		Q: req.GetQ(),
	}
	
	return reqVO
}

func (c CandidateListHandler) BuildResponse(resVO vo.CandidateListResVO) []*irisProtoc.Candidate  {
	var (
		response           []*irisProtoc.Candidate
	)

	candidates := resVO.Candidates
	if len(candidates) > 0 {
		for _, v := range candidates {
			var (
				resCandidate irisProtoc.Candidate
				resCandidateDescription irisProtoc.CandidateDescription
				resCandidateDelegator   irisProtoc.Delegator
				
				resCandidateDelegators []*irisProtoc.Delegator
			)
			// description
			resCandidateDescription = irisProtoc.CandidateDescription{
				Details: v.Description.Details,
				Identity: v.Description.Identity,
				Moniker: v.Description.Moniker,
				Website: v.Description.Website,
			}
			
			// delegator
			if len(v.Delegators) > 0 {
				delegator := v.Delegators[0]
				resCandidateDelegator = irisProtoc.Delegator{
					Address: delegator.Address,
					PubKey: delegator.ValidatorAddr,
					Shares: delegator.Shares,
				}
				resCandidateDelegators = append(resCandidateDelegators, &resCandidateDelegator)
			}
			
			
			resCandidate = irisProtoc.Candidate{
				Address: v.Address,
				PubKey: v.PubKey,
				Shares: v.Shares,
				VotingPower: v.VotingPower,
				Description: &resCandidateDescription,
				Delegators: resCandidateDelegators,
			}
			
			response = append(response, &resCandidate)
		}
	}
	return response
}
