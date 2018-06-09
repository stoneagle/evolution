import { Component, OnInit, ViewChild } from '@angular/core';
import { Pool } from '../../model/business/pool';
import { PoolService } from '../../service/business/pool.service';
import { SavePoolComponent } from './save/save.component';

@Component({
  selector: 'app-pool',
  templateUrl: './pool.component.html',
  styleUrls: ['./pool.component.css']
})
export class PoolComponent implements OnInit {
  @ViewChild(SavePoolComponent)
  savePool: SavePoolComponent;

  pools: Pool[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private poolService: PoolService
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  saved(saved: boolean): void {
    if (saved) {
      this.refresh();
    }
  }

  openSaveModel(id?: number): void {
    this.savePool.New(id);
  }

  delete(pool: Pool): void {
    this.poolService.Delete(pool.Id)
    .subscribe(autobuild => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshPools(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshPools(0, 10);
  }

  refreshPools(from: number, to: number): void {
    this.poolService.List()
    .subscribe(res => {
      this.totalCount = res.length;
      this.pools = res.slice(from, to);
    })
  }
}
