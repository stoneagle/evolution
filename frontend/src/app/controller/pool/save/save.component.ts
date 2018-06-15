import { Component, OnInit, Output, EventEmitter } from '@angular/core'; 
import { Pool } from '../../../model/business/pool';
import { PoolService  } from '../../../service/business/pool.service';

@Component({
  selector: 'save-pool',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class SavePoolComponent implements OnInit {
  pool: Pool = new Pool;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private poolService: PoolService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.poolService.Get(id)
      .subscribe(res => {
        this.pool = res;
        this.modelOpened = true;
      })
    } else {
      this.pool = new Pool();
      this.modelOpened = true;
    }
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
}
