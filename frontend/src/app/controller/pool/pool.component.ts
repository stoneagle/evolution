import { Component, OnInit } from '@angular/core';
import { Pool } from '../../model/pool';
import { PoolService } from '../../service/business/pool.service';

@Component({
  providers: [PoolService],
  selector: 'app-pool',
  templateUrl: './pool.component.html',
  styleUrls: ['./pool.component.css']
})
export class PoolComponent implements OnInit {
  pools: Pool[];

  constructor(
    private poolService: PoolService
  ) { }

  ngOnInit() {
    this.getPools();
  }

  getPools(): void {
    this.poolService.getPools()
    .subscribe(pools => this.pools = pools);
  }
}
