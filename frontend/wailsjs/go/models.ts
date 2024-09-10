export namespace models {
	
	export class LatestTransaction {
	    user_id: string;
	    user_name: string;
	    // Go type: time
	    transaction_date: any;
	    cumulative_transaction_count: number;
	
	    static createFrom(source: any = {}) {
	        return new LatestTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.transaction_date = this.convertValues(source["transaction_date"], null);
	        this.cumulative_transaction_count = source["cumulative_transaction_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TransactionLeaderboard {
	    user_id: string;
	    user_name: string;
	    current_rank: number;
	    rank_change_indicator: string;
	    total_transaction_count: number;
	
	    static createFrom(source: any = {}) {
	        return new TransactionLeaderboard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.current_rank = source["current_rank"];
	        this.rank_change_indicator = source["rank_change_indicator"];
	        this.total_transaction_count = source["total_transaction_count"];
	    }
	}

}

