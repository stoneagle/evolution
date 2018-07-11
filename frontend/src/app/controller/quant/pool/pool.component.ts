import { Component, OnInit, ViewChild } from '@angular/core';
import { Router, ActivatedRoute }       from "@angular/router";
import { Pool }                         from '../../../model/quant/pool';
import { PoolService }                  from '../../../service/quant/pool.service';
import { PoolSaveComponent }            from './save/save.component';

@Component({
  selector: 'quant-pool',
  templateUrl: './pool.component.html',
  styleUrls: ['./pool.component.css']
})
export class PoolComponent implements OnInit {
  @ViewChild(PoolSaveComponent)
  saveComponent: PoolSaveComponent;

  pools: Pool[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private poolService: PoolService,
    private router: Router, 
    private route: ActivatedRoute 
  ) { 
  }

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
    switch(this.route.parent.snapshot.routeConfig.path) {
      case 'stock':
        this.saveComponent.New("Stock", id);
        break;
    }
  }

  delete(pool: Pool): void {
    this.poolService.Delete(pool.Id)
    .subscribe(res => {
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
    this.poolService.List().subscribe(res => {
      this.totalCount = res.length;
      this.pools = res.slice(from, to);
    })
  }

  enterItemList(id: number):void {
    let linkUrl = ['stock', 'pool', id]; 
    this.router.navigate(linkUrl);
  }
}
