package query

import "github.com/CarosDrean/api-results.git/models"

var protocolcomponent = models.TableDB{
	Name:   "dbo.protocolcomponent",
	Fields: []string{"v_ProtocolComponentId", "v_ProtocolId", "r_Price"},
}

var ProtocolComponent = models.QueryDB{
	"listComponent": {Q: "select co.v_ComponentId, co.v_Name,c.r_Price from protocol p" +
		" inner join protocolcomponent c on c.v_ProtocolId  = p.v_ProtocolId" +
		" inner join component co on c.v_ComponentId = co.v_ComponentId " +
		" where p.v_ProtocolId = '%s' and p.i_IsDeleted = 0 and co.i_IsDeleted = 0 " +
		" group by co.v_Name, co.v_ComponentId, c.r_Price"},
}

