import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Pool }           from '../../../../model/quant/pool';
import { PoolService  }   from '../../../../service/quant/pool.service';
import { ConfigService  } from '../../../../service/quant/config.service';

@Component({
  selector: 'quant-pool-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class PoolSaveComponent implements OnInit {
  typeMap: Map<string, string> = new Map();
  strategyMap: Map<string, string> = new Map();
  pool: Pool = new Pool;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private poolService: PoolService,
    private configService: ConfigService,
  ) { }

  ngOnInit() {
  }

  New(asset:string, id?: number): void {
    this.configService.TypeList(asset)
      .subscribe(res => {
        res.forEach((mainMap, ctype) => {
          this.typeMap.set(ctype, ctype);
        });
        if (id) {
          this.poolService.Get(id)
          .subscribe(res => {
            this.pool = res;
            this.modelOpened = true;
          })
        } else {
          this.pool = new Pool();
          this.pool.AssetString = asset;
          this.modelOpened = true;
        }
      })
  }            

  TypeOnChange(ctype) {
    this.configService.StrategyList(ctype).subscribe(res => {
      this.strategyMap = res;
      this.pool.Type = ctype; 
      this.pool.Strategy = this.strategyMap.keys().next().value;
    });
  }

  Submit(): void {
    if (this.pool.Id == null) {
      this.poolService.Add(this.pool)
      .subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.poolService.Update(this.pool)
      .subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }

  getKeys(map) {
    return Array.from(map.keys());
  }
}
